package rest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_AppInfo(t *testing.T) {
	router := chi.NewRouter()
	router.With(AppInfo("app-name", "Umputun", "12345")).Get("/blah", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("blah blah"))
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/blah")
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "blah blah", string(b))
	assert.Equal(t, "app-name", resp.Header.Get("App-Name"))
	assert.Equal(t, "12345", resp.Header.Get("App-Version"))
	assert.Equal(t, "Umputun", resp.Header.Get("Org"))
}

func TestMiddleware_GetBodyAndUser(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/request", strings.NewReader("body"))
	require.Nil(t, err)

	body, user := getBodyAndUser(req, []LoggerFlag{LogAll})
	assert.Equal(t, "body", body)
	assert.Equal(t, "", user, "no user")

	req = SetUserInfo(req, uinfo{id: "id1", name: "user1"})
	_, user = getBodyAndUser(req, []LoggerFlag{LogAll})
	assert.Equal(t, ` - user1/id1`, user, "no user")

	body, user = getBodyAndUser(req, nil)
	assert.Equal(t, "", body)
	assert.Equal(t, "", user, "no user")

	body, user = getBodyAndUser(req, []LoggerFlag{LogNone})
	assert.Equal(t, "", body)
	assert.Equal(t, "", user, "no user")

	body, user = getBodyAndUser(req, []LoggerFlag{LogUser})
	assert.Equal(t, "", body)
	assert.Equal(t, ` - user1/id1`, user, "no user")
}