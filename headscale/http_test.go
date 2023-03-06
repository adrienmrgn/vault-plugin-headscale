package headscale

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"fmt"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestBuildHTTPRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := newClient()
	c.ApiURL = ts.URL
	c.ApiKey = "foobarbaz"

	// var req *http.Request
	reqOptWithNil := requestOptions{
		context: context.Background(),
		uri: "/",
		method: httpGet,
		queryParams: nil,
		queryBody: nil,
	}
	req, err := c.buildHTTPRequest(reqOptWithNil)
	assert.NoError(t, err, "No error during building request with nil body and paramas")
	assert.Nil(t, err)

	reqOptWithParams := requestOptions{
		context: context.Background(),
		uri: "/",
		method: httpGet,
		queryParams: map[string]string{
			"user": "foo",
		},
		queryBody: nil,
	}
	req, err = c.buildHTTPRequest(reqOptWithParams)
	assert.NoError(t, err, "No error during building request with query param")
	assert.Nil(t, err)

	reqOptWithBody := requestOptions{
		context: context.Background(),
		uri: "/",
		method: httpGet,
		queryParams: nil,
		queryBody: map[string]any{
			"user":       "foo",
			"used":       true,
			"reusable":   false,
			"ephemeral":  true,
			"expiration": time.Now(),
			"acl_tags":   []string{"tag1","tag2"},
		},
	}
	req, err = c.buildHTTPRequest(reqOptWithBody)
	assert.NoError(t, err, "No error during building request with query param")
	assert.Nil(t, err)

	_, err = c.HTTP.Do(req)
	assert.NoError(t, err, "No error during executing the request")

}

func TestGet(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query, _ := url.ParseQuery(r.URL.RawQuery)
		param := query.Get("user")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "GET query param user : %s", param)

	}))
	defer ts.Close()

	c := newClient()
	c.ApiURL = ts.URL
	c.ApiKey = "foobarbaz"
	queryParam := map[string]string{
		"user": "foo",
	}
	resp, err := c.get(context.Background(),"/",queryParam)
	closeResponseBody(resp)
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response if %s",data)
}