# httputil
http util for restful request: get, post, put, patch, delete

## Options (WithFunc)

http request options

### WithClient

default: `&http.Client{}`

example: client with https

```golang
tr := &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
client := &http.Client{Transport: tr}
```

### WithTimeout

default: `30 * time.Second`

### WithHeader

default: `"Accept": "application/json", "Content-Type": "application/json;charset=UTF-8"`

### WithMarshal

default: `json.Marshal`

### WithUnmarshal

default: `json.Unmarshal`

### WithLogTimeCost

default: donot log request time cost, using logger.Printf

log style: gin

example: `INFO 2021/05/04 12:24:04 REQUEST | 200 |  305.549239ms | GET     https://httpbin.org/get?hello=world`

### WithStatusCodeJudge

default: `defaultCodeJudger` (2xx is the right status code)

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