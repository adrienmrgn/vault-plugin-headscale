package headscale

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildHTTPRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient()
	c.APIURL = ts.URL
	c.APIKey = "foobarbaz"

	// var req *http.Request
	reqOptWithNil := requestOptions{
		context:     context.Background(),
		uri:         "/",
		method:      httpGet,
		queryParams: nil,
		queryBody:   nil,
	}
	req, err := c.buildHTTPRequest(reqOptWithNil)
	assert.NoError(t, err, "No error during building request with nil body and paramas")
	assert.Nil(t, err)

	reqOptWithParams := requestOptions{
		context: context.Background(),
		uri:     "/",
		method:  httpGet,
		queryParams: map[string]string{
			"user": "foo",
		},
		queryBody: nil,
	}
	req, err = c.buildHTTPRequest(reqOptWithParams)
	assert.NoError(t, err, "No error during building request with query param")
	assert.Nil(t, err)

	reqOptWithBody := requestOptions{
		context:     context.Background(),
		uri:         "/",
		method:      httpGet,
		queryParams: nil,
		queryBody: map[string]any{
			"user":       "foo",
			"used":       true,
			"reusable":   false,
			"ephemeral":  true,
			"expiration": time.Now(),
			"acl_tags":   []string{"tag1", "tag2"},
		},
	}
	req, err = c.buildHTTPRequest(reqOptWithBody)
	assert.NoError(t, err, "No error during building request with query param")
	assert.Nil(t, err)

	_, err = c.HTTP.Do(req)
	assert.NoError(t, err, "No error during executing the request")

}

func TestGet(t *testing.T) {
	responseTemplate := "GET query param user : %s"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}
		query, _ := url.ParseQuery(r.URL.RawQuery)
		param := query.Get("user")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, responseTemplate, param)

	}))
	defer testServer.Close()

	c := NewClient()
	c.APIURL = testServer.URL
	c.APIKey = "foobarbaz"

	queryParam := map[string]string{
		"user": "foo",
	}
	resp, err := c.get(context.Background(), "/", queryParam)
	defer closeResponseBody(resp)

	assert.NoError(t, err)
	data, err := io.ReadAll(resp.Body)
	expectedREsponse := fmt.Sprintf(responseTemplate, queryParam["user"])
	assert.NoError(t, err)
	assert.Equal(t, expectedREsponse, string(data), "HTTP response matches")
}

func TestPost(t *testing.T) {

	userName := "foo"
	expected := "Success"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifie que la méthode est bien POST
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}

		// Vérifie le type de contenu
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}

		// Vérifie le corps de la requête
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Vérifie les champs de la requête
		if data["name"] != userName {
			http.Error(w, "Invalid request fields", http.StatusBadRequest)
			return
		}

		// Envoie une réponse réussie
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))
	defer testServer.Close()
	c := NewClient()
	c.APIURL = testServer.URL
	c.APIKey = "foobarbaz"
	requestBody := map[string]string{
		"name": userName,
	}
	resp, err := c.post(context.Background(), "/", requestBody)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	var actualBody bytes.Buffer
	_, err = actualBody.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	actual := actualBody.String()
	assert.Equal(t, expected, actual)
}
