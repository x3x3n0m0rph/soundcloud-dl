<p align="center">
    <img src="assets/go-sound.png" height="400px" width="100%" style="border-radius:8px" alt="goandsoundcloud">
</p>

<h1 align="center">
SoundCloud-dl

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Sound Cloud](https://img.shields.io/badge/sound%20cloud-FF5500?style=for-the-badge&logo=soundcloud&logoColor=white)

</h1>

A simple CLI app written in GO used to download sound-tracks from soundcloud with a nice UI and some cool features like downloading with multiple formats/qualities and search.

# Overview

`soundcloud-dl` is a compiled CLI app that uses the soundcloud API to :

- Download tracks by URL.
- Search in-app and download.
- Save tracks with tags/metadata : `Artwork Image`, `Creation Date`, `Author Name`.
- Multiple formats and qualities.
- [Blazingly Fast](https://youtu.be/Z0GX2mTUtfo)
- Download playlists.
- Download user's tracks and playlists. (coming soon).

# Installation

There are multiple ways to install, the easiest is :
```
go install github.com/x3x3n0m0rph/soundcloud-dl@latest
```
other way is to grab is the source code and build.

# How to use ?

Make sure to read the help menu.

```
Usage:
  soundclound-dl <url> [flags]

Flags:
  -b, --best                   Download with the best available quality.
  -f, --force                  Overwrite existing files.
  -p, --download-path string   The download path where tracks are stored. (default "/home/none/Things/github/soundcloud-dl")
  -h, --help                   help for sc
  -s, --search-and-download    Search for tracks by title and prompt one for download 
      --socks-proxy string     SOCKS5 proxy connection string, for example socks5://user:pass@host:1080.
  -v, --version                version for sc
```

Notes : `-s` can't work with a `url` passed.

SOCKS5 proxy can be passed as a connection string:

```
sc <url> --socks-proxy "socks5://user:pass@host:8888"
```

- **Download through a given URL**
<p align="center">
    <img src="assets/soundcloud-dl-github.gif" style="border-radius:8px" alt="goandsoundcloud">
</p>

- **Download through search**

Note: search is limited to 6 results of type `track`, later on maybe `playlist` will be added. Also pagination may be add!.

<p align="center">
    <img src="assets/download-search-short.gif" style="border-radius:8px" alt="goandsoundcloud">
</p>

- **Download a playlist**

Note : when downloading a playlist, it uses the flag `--best` to grab the best quality.

<p align="center">
    <img src="assets/playlist-download.gif" style="border-radius:8px" alt="goandsoundcloud">
</p>
