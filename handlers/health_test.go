package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	asserter := assert.New(t)

	r := chi.NewRouter()
	r.Get("/health", Health())
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/health", ts.URL), nil)
	res, err := ts.Client().Do(req)
	asserter.NoError(err)
	defer res.Body.Close()

	asserter.Equal(http.StatusOK, res.StatusCode)
}
