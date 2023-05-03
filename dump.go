package httputil

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func CurlLikeDumpRequest(req *http.Request) {
	data, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return
	}
	printCurlLikeDump(data, "> ")
}

func CurlLikeDumpResponse(resp *http.Response) {
	data, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return
	}
	printCurlLikeDump(data, "< ")
}

func printCurlLikeDump(data []byte, lineStart string) {
	start := fmt.Sprintf("\n%s", lineStart)
	tmp := bytes.ReplaceAll(data, []byte("\n"), []byte(start))
	tmp = bytes.TrimRight(tmp, lineStart)
	fmt.Printf("%s%s", lineStart, tmp)
}
