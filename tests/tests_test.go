package tests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestFunctional(t *testing.T) {
	githubPushFile, err := os.Open("testdata/test-repo-generated.json")
	if err != nil {
		t.Fatal(err)
	}
	defer githubPushFile.Close()
	dockerhubFile, err := os.Open("testdata/dockerhub-push.json")
	if err != nil {
		t.Fatal(err)
	}
	defer dockerhubFile.Close()
	hubSignature, err := ioutil.ReadFile("testdata/test-repo-generated.hash")
	if err != nil {
		t.Fatal(err)
	}
	requests := []*http.Request{
		{
			Method: "POST",
			URL:    &url.URL{Scheme: "http", Host: "localhost:7744", Path: "/events/github"},
			Body:   githubPushFile,
			Header: http.Header{
				"X-Github-Event":  []string{"push"},
				"X-Hub-Signature": []string{string(hubSignature)},
			},
		},
		{
			Method: "POST",
			URL:    &url.URL{Scheme: "http", Host: "localhost:7744", Path: "/events/dockerhub/deis/empty-testbed/589e15029e1e44dee48de4800daf1f78e64287c0"},
			Body:   dockerhubFile,
		},
	}

	for _, request := range requests {
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s %s: expected status code '200', got '%d'\n", request.Method, request.URL.String(), resp.StatusCode)
		}
	}
}
