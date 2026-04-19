package client

import (
	"net/http"
	"net/url"
	"testing"
)

func TestConfigureSocksProxy(t *testing.T) {
	t.Cleanup(func() {
		if err := Configure(""); err != nil {
			t.Fatalf("failed to reset client: %s", err)
		}
	})

	proxyURL := "socks5://user:pass@localhost:1080"
	if err := Configure(proxyURL); err != nil {
		t.Fatalf("expected SOCKS5 proxy URL to be accepted: %s", err)
	}

	transport, ok := activeClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected configured transport, got %T", activeClient.Transport)
	}

	reqURL, err := url.Parse("https://soundcloud.com/test")
	if err != nil {
		t.Fatal(err)
	}

	got, err := transport.Proxy(&http.Request{URL: reqURL})
	if err != nil {
		t.Fatalf("proxy lookup failed: %s", err)
	}
	if got.String() != proxyURL {
		t.Fatalf("expected proxy %q, got %q", proxyURL, got.String())
	}
}
