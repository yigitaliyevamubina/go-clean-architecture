package integrationtest

import (
	"log"
	"net/http"
	"os"
	"testing"

	. "github.com/Eun/go-hit"
)

const (
	// Attempts connection
	host       = "localhost:8080"
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath = "/http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", host, attempts)

		attempts--
	}

	return err
}

// HTTP POST: /post/create.
func TestHTTPCreatePost(t *testing.T) {
	body := `{
		"user_id": "d0b69f3b-2021-4d91-8e13-c243d9eb5292",
		"content": "This is the content of post 13.",
		"title": "Post 13",
		"likes": 20,
		"dislikes": 6,
		"views": 100,
		"category": "Nature"
	}`
	Test(t, 
		Description("Create post Success"),
		Post(basePath+"/post/create"), 
		Send().Headers("Content-Type").Add("application/json"), 
		Send().Body().String(body), 
		Expect().Status().Equal(http.StatusOK), 
		Expect().Body().JSON().JQ(".title").Equal("Post 13"),	
	)

	body = `{
		"user_id": "",
		"content": "This is the content of post 13.",
		"title": "Post 13",
		"likes": 20,
		"dislikes": 6,
		"views": 100,
		"category": "Nature"
	}`
	Test(t, 
		Description("Create post Fail"),
		Post(basePath+"/post/create"), 
		Send().Headers("Content-Type").Add("application/json"), 
		Send().Body().String(body), 
		Expect().Status().Equal(http.StatusBadRequest), 
		Expect().Body().JSON().JQ(".error").Equal("invalid request body"),	
	)
}

// HTTP LIST: /posts/:page/:limit
func TestHTTPListPosts(t *testing.T) {
	Test(t,
		Description("History Success"), 
		Get(basePath+"/posts/1/10"), 
		Expect().Status().Equal(http.StatusOK), 
		Expect().Body().String().Contains(`{"posts":[{`),
	)
}