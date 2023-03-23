package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerWorks(t *testing.T) {
	t.Run("an arbitrary route returns something", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/req", nil)
		res := httptest.NewRecorder()

		URLServer(res, req)

		got := res.Body.String()
		want := "works"

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}
