# httputil
http util for restful request: get, post, put, patch, delete

## Usage

```golang
package main

import (
	"context"
	"log"

	hu "github.com/chinaran/httputil"
)

func main() {
	// get
	urlGet := "https://httpbin.org/get?hello=world"
	respGetM := map[string]interface{}{}
	if err := hu.Get(context.TODO(), urlGet, &respGetM, hu.WithLogTimeCost()); err != nil {
		log.Printf("Get %s err: %s", urlGet, err)
		return
	}
	log.Printf("Get %s map response: %+v", urlGet, respGetM)
	respGetStr := ""
	if err := hu.Get(context.TODO(), urlGet, &respGetStr, hu.WithLogTimeCost()); err != nil {
		log.Printf("Get %s err: %s", urlGet, err)
		return
	}
	log.Printf("Get %s string response: %+v", urlGet, respGetStr)

	// post
	urlPost := "https://httpbin.org/post"
	req := map[string]string{"hello": "world"}
	respPost := struct {
		Data string `json:"data"`
	}{}
	if err := hu.Post(context.TODO(), urlPost, &req, &respPost, hu.WithLogTimeCost()); err != nil {
		log.Printf("Post %s err: %s", urlGet, err)
		return
	}
	log.Printf("Post %s struct response: %+v", urlGet, respPost)
}
```