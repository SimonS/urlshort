package urlshort_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshort"
)

func stubURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "works")
}

func StubMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", stubURLHandler)
	return mux
}

func TestServerWorks(t *testing.T) {
	stub := StubMux()

	t.Run("an arbitrary route returns something", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		stub.ServeHTTP(res, req)

		got := res.Body.String()
		want := "works"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})

	t.Run("map handler takes a dictionary and redirects appropriately", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/redirect", nil)
		res := httptest.NewRecorder()

		pathsToUrls := map[string]string{
			"/redirect": "https://mytotallycoolwebsite.com",
		}
		mappedServer := urlshort.MapHandler(pathsToUrls, stub)
		mappedServer(res, req)

		assertStatusCode(t, res, http.StatusPermanentRedirect)

		assertLocation(t, res, "https://mytotallycoolwebsite.com")
	})

	t.Run("map handler 404s on a URL it doesn't know about", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/404", nil)
		res := httptest.NewRecorder()

		pathsToUrls := map[string]string{
			"/redirect": "https://mytotallycoolwebsite.com",
		}
		mappedServer := urlshort.MapHandler(pathsToUrls, stub)
		mappedServer(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("map handler takes a fallback handler", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/404", nil)
		res := httptest.NewRecorder()

		pathsToUrls := map[string]string{
			"/redirect": "https://mytotallycoolwebsite.com",
		}
		mappedServer := urlshort.MapHandler(pathsToUrls, stub)
		mappedServer(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})
}

func assertStatusCode(t testing.TB, r *httptest.ResponseRecorder, want int) {
	t.Helper()

	got := r.Result().StatusCode

	if got != want {
		t.Errorf("got status %d, but wanted %d", got, want)
	}
}

func assertLocation(t testing.TB, r *httptest.ResponseRecorder, want string) {
	t.Helper()

	got := r.Result().Header.Get("Location")

	if got != want {
		t.Errorf("got location %q, but wanted %q", got, want)
	}
}
