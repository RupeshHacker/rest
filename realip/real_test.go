package realip

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFromHeaders(t *testing.T) {
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		req.Header.Add("Something", "1234567")
		req.Header.Add("X-Real-IP", "8.8.8.8")
		adr, err := Get(req)
		require.NoError(t, err)
		assert.Equal(t, "8.8.8.8", adr)
	}
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		req.Header.Add("Something", "1234567")
		req.Header.Add("X-Forwarded-For", "8.8.8.8,1.1.1.2, 30.30.30.1")
		adr, err := Get(req)
		require.NoError(t, err)
		assert.Equal(t, "30.30.30.1", adr)
	}
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		req.Header.Add("Something", "1234567")
		req.Header.Add("X-Forwarded-For", "8.8.8.8,1.1.1.2,192.168.1.1,10.0.0.65")
		adr, err := Get(req)
		require.NoError(t, err)
		assert.Equal(t, "1.1.1.2", adr)
	}
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		req.Header.Add("Something", "1234567")
		req.Header.Add("X-Forwarded-For", "30.30.30.1")
		req.Header.Add("X-Real-Ip", "10.0.0.1")
		adr, err := Get(req)
		require.NoError(t, err)
		assert.Equal(t, "30.30.30.1", adr)
	}
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		req.Header.Add("Something", "1234567")
		req.Header.Add("X-Forwarded-For", "30.30.30.1")
		req.Header.Add("X-Real-Ip", "8.8.8.8")
		adr, err := Get(req)
		require.NoError(t, err)
		assert.Equal(t, "30.30.30.1", adr)
	}
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		req.Header.Add("Something", "1234567")
		req.Header.Add("X-Forwarded-For", "10.0.0.2,192.168.1.1")
		req.Header.Add("X-Real-Ip", "8.8.8.8")
		adr, err := Get(req)
		require.NoError(t, err)
		assert.Equal(t, "8.8.8.8", adr)
	}
	{
		req, err := http.NewRequest("GET", "/something", http.NoBody)
		assert.NoError(t, err)
		ip, err := Get(req)
		assert.Error(t, err)
		assert.Equal(t, "", ip)
	}
}

func TestGetFromRemoteAddr(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r)
		adr, err := Get(r)
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1", adr)
	}))

	req, err := http.NewRequest("GET", ts.URL+"/something", http.NoBody)
	require.NoError(t, err)
	client := http.Client{Timeout: time.Second}
	_, err = client.Do(req)
	require.NoError(t, err)
}