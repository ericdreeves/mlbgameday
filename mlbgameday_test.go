package mlbgameday

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(nil)

	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if r.Method != want {
		t.Errorf("r.Method == %v, want %v", r.Method, want)
	}
}

func testError(t *testing.T, subject string, want string, err error) {
	if err == nil {
		t.Fatalf("%v did not return an error", subject)
	}

	if err.Error() != want {
		t.Errorf("%v returned error %q, want %q", subject, err.Error(), want)
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	client := NewClient(nil)
	got := client.BaseURL.String()
	want := "http://gd2.mlb.com"
	if got != want {
		t.Errorf("client.BaseURL.String() == %q, want %q", got, want)
	}
}
