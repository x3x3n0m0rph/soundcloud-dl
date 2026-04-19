package soundcloud_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/grafov/m3u8"
)

func TestDownload(t *testing.T) {

	fileResp := readTestFile("Test[low].mp3")
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write(fileResp)
	}))

	defer testServer.Close()

	path := "../../testdata/"
	downloadTrack := soundcloud.DownloadTrack{
		Url:     testServer.URL,
		Quality: "low",
		Ext:     "mp3",
		SoundData: &soundcloud.SoundData{
			Id:           75825848,
			CreatedAt:    "2013-01-21T13:14:46Z",
			PermalinkUrl: "https://soundcloud.com/sobhi-mohamed5/99-118-mp4",
			ArtworkUrl:   "https://i1.sndcdn.com/artworks-000038806208-j0yb4i-large.jpg",
			Title:        "Test1",
			LikesCount:   355,
		},
	}

	expectedPath, err := soundcloud.Download(downloadTrack, path, false)
	if err != nil {
		t.Fatalf("expected download to succeed: %s", err)
	}

	// read the downloaded file
	file, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		t.Errorf("An error happened while reading the track, error : %s", err)
	}
	// TODO: not the best method, as it loads all the file in memeory, but for this test it's ok since the size isn't that big + I have RAM
	if !bytes.Equal(file, fileResp) {
		t.Errorf("Expected the 2 files to be the same")
	}

	// remove the file
	os.Remove(expectedPath)
}

func TestDownloadForceOverwritesExistingFile(t *testing.T) {
	fileResp := []byte("new content")
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write(fileResp)
	}))
	defer testServer.Close()

	path := t.TempDir()
	downloadTrack := soundcloud.DownloadTrack{
		Url:     testServer.URL,
		Quality: "low",
		Ext:     "mp3",
		SoundData: &soundcloud.SoundData{
			Title: "TestForce",
		},
	}

	expectedPath := filepath.Join(path, "TestForce[low].mp3")
	if err := os.WriteFile(expectedPath, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	gotPath, err := soundcloud.Download(downloadTrack, path, false)
	if !errors.Is(err, soundcloud.ErrFileExists) {
		t.Fatalf("expected existing file to be skipped without force, got path %s and error %v", gotPath, err)
	}
	if gotPath != "" {
		t.Fatalf("expected existing file to be skipped without force, got %s", gotPath)
	}

	gotPath, err = soundcloud.Download(downloadTrack, path, true)
	if err != nil {
		t.Fatalf("expected force download to succeed: %s", err)
	}
	if gotPath != expectedPath {
		t.Fatalf("expected overwritten path %s, got %s", expectedPath, gotPath)
	}

	file, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(file, fileResp) {
		t.Fatalf("expected file to be overwritten")
	}
}

func extractSegments(fileP []byte, testfile []byte) map[int][]byte {
	segments := make(map[int][]byte, 0)
	reader := bytes.NewReader(fileP)

	pl, listType, err := m3u8.DecodeFrom(reader, true)

	if err != nil {
		return nil
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := pl.(*m3u8.MediaPlaylist)
		for i, segment := range mediapl.Segments {
			if segment == nil {
				continue
			}
			segments[i] = nil
		}
	}
	var segmentSize int
	numSegments := len(segments)

	fileSize := len(testfile)
	segmentSize = fileSize / numSegments

	for i := 0; i < numSegments-1; i++ {
		start := i * segmentSize
		end := (i + 1) * segmentSize
		segments[i] = testfile[start:end]
	}

	// handle the last segment
	start := (numSegments - 1) * segmentSize
	end := fileSize
	if fileSize%numSegments == 0 {
		end = (numSegments-1)*segmentSize + segmentSize
	}
	segments[int(numSegments-1)] = testfile[start:end]

	return segments
}

func TestDownloadM3u8(t *testing.T) {

	path := "../../testdata/"
	downloadTrack := soundcloud.DownloadTrack{
		Url:     "",
		Quality: "medium",
		Ext:     "mp3",
		SoundData: &soundcloud.SoundData{
			Id:           75825848,
			CreatedAt:    "2013-01-21T13:14:46Z",
			PermalinkUrl: "https://soundcloud.com/sobhi-mohamed5/99-118-mp4",
			ArtworkUrl:   "https://i1.sndcdn.com/artworks-000038806208-j0yb4i-large.jpg",
			Title:        "Test1",
			LikesCount:   355,
		},
	}
	fileResp := readTestFile("playlist.m3u8")
	track := readTestFile("Test[medium].mp3")

	fileName := downloadTrack.SoundData.Title + "[" + downloadTrack.Quality + "]." + downloadTrack.Ext
	path = filepath.Join(path, fileName)
	segments := extractSegments(fileResp, track)

	// modifiying the seg urls to mimic the actual urls
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// write bytes based on the url
		n, _ := strconv.Atoi(req.URL.String()[1:])
		if bs, ok := segments[n]; ok {
			res.WriteHeader(http.StatusOK)
			res.Write(bs)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write(fileResp)
	}))

	// setting the server url
	downloadTrack.Url = testServer.URL

	defer testServer.Close()

	segmentURIs := make([]string, 0)
	keys := make([]int, 0, len(segments))
	for k := range segments {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		segmentURIs = append(segmentURIs, testServer.URL+"/"+strconv.Itoa(k))
	}
	soundcloud.DownloadM3u8(path, nil, segmentURIs)

	// read the downloaded file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("An error happened while reading the track, error : %s", err)
	}
	// TODO: not the best method, as it loads all the file in memeory, but for this test it's ok since the size isn't that big + I have RAM
	if !bytes.Equal(file, track) {
		t.Errorf("Expected the 2 files to be the same")
	}
	// remove the file
	os.Remove(path)
}
