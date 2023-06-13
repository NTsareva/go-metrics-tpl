package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return resp, string(respBody)
}

func TestPostRouter(t *testing.T) {
	ts := httptest.NewServer(MetricsRouter())
	defer ts.Close()

	var testTable = []struct {
		url    string
		status int
	}{
		{"/update/gauge/alloc/1111.111", http.StatusOK},
		{"/update/counter/myvalue/1", http.StatusOK},
		// проверим на ошибочный запрос
		{"/update", http.StatusBadRequest},
		{"/update/", http.StatusBadRequest},
		{"/update/gauge", http.StatusNotFound},
		{"/update/gauge/", http.StatusNotFound},
		{"/update/counter/", http.StatusNotFound},
		{"/update/counter", http.StatusNotFound},
		{"/update/gauge/alloc/", http.StatusBadRequest},
		{"/update/gauge/alloc", http.StatusBadRequest},
		{"/update/counter/mymetrics/", http.StatusBadRequest},
		{"/update/counter/mymetrics", http.StatusBadRequest},
	}
	for _, v := range testTable {
		resp, _ := testRequest(t, ts, "POST", v.url)
		assert.Equal(t, v.status, resp.StatusCode)
	}
}
