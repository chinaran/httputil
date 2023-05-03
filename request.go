package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// default logger
var logger = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)

// PrintfFunc for print request log
type PrintfFunc func(format string, v ...interface{})

// StatusCodeJudgeFunc for judge status code (right status code return true)
type StatusCodeJudgeFunc func(statusCode int) bool

// MarshalFunc marshal request function
type MarshalFunc func(v interface{}) ([]byte, error)

// UnmarshalFunc unmarshal response function
type UnmarshalFunc func(data []byte, v interface{}) error

// request options
type reqOptions struct {
	httpClient   *http.Client
	timeout      time.Duration
	header       map[string]string
	marshaler    MarshalFunc
	unmarshaler  UnmarshalFunc
	logTimeCost  bool
	dumpRequest  bool
	dumpResponse bool
	printfer     PrintfFunc
	codeJudger   StatusCodeJudgeFunc
	statusCode   *int
}

// ReqOptionFunc request option function
type ReqOptionFunc func(opt *reqOptions) error

// WithClient default: &http.Client{}
func WithClient(client *http.Client) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil || client == nil {
			return nil
		}
		opt.httpClient = client
		return nil
	}
}

// WithTimeout default: 30 * time.Second
func WithTimeout(timeout time.Duration) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil {
			return nil
		}
		opt.timeout = timeout
		return nil
	}
}

// WithHeader default: "Accept": "application/json", "Content-Type": "application/json;charset=UTF-8"
func WithHeader(header map[string]string) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil || header == nil {
			return nil
		}
		opt.header = header
		return nil
	}
}

// WithTraceHeaders append trace headers
// NOTE:
// 1. This should only be used in the demo project.
// 2. If the WithHeader method is used, it should be used after it.
func WithTraceHeaders(headers http.Header) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil || headers == nil {
			return nil
		}
		if opt.header == nil {
			opt.header = make(map[string]string)
		}
		m := ForwardHeaders(headers)
		for k, v := range m {
			opt.header[k] = v
		}
		return nil
	}
}

// WithMarshal default: json.Marshal
func WithMarshal(marshaler MarshalFunc) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil || marshaler == nil {
			return nil
		}
		opt.marshaler = marshaler
		return nil
	}
}

// WithUnmarshal default: json.Unmarshal
func WithUnmarshal(unmarshaler UnmarshalFunc) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil || unmarshaler == nil {
			return nil
		}
		opt.unmarshaler = unmarshaler
		return nil
	}
}

// WithLogTimeCost default: false, logger.Printf
func WithLogTimeCost(printfer ...PrintfFunc) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil {
			return nil
		}
		opt.logTimeCost = true
		if len(printfer) > 0 {
			opt.printfer = printfer[0]
		}
		return nil
	}
}

// WithDumpRequest default: false, false
func WithDumpRequest(dumpRequest, dumpResponse bool) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil {
			return nil
		}
		opt.dumpRequest = dumpRequest
		opt.dumpResponse = dumpResponse
		return nil
	}
}

// WithStatusCodeJudge default: defaultCodeJudger
func WithStatusCodeJudge(codeJudger StatusCodeJudgeFunc) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil || codeJudger == nil {
			return nil
		}
		opt.codeJudger = codeJudger
		return nil
	}
}

// default StatusCodeJudgeFunc
func defaultCodeJudger(statusCode int) bool {
	if statusCode < 200 || statusCode >= 300 {
		return false
	}
	return true
}

// StoreStatusCode default: not store, just judge
func StoreStatusCode(statusCode *int) ReqOptionFunc {
	return func(opt *reqOptions) error {
		if opt == nil {
			return nil
		}
		opt.statusCode = statusCode
		return nil
	}
}

// default option for each request
var defaultReqOption = reqOptions{
	httpClient: &http.Client{},
	timeout:    30 * time.Second,
	header: map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json;charset=UTF-8",
	},
	marshaler:    json.Marshal,
	unmarshaler:  json.Unmarshal,
	logTimeCost:  false,
	dumpRequest:  false,
	dumpResponse: false,
	printfer:     logger.Printf,
	codeJudger:   defaultCodeJudger,
	statusCode:   new(int),
}

// Get http request
func Get(ctx context.Context, url string, resp interface{}, opts ...ReqOptionFunc) error {
	return HttpRequest(ctx, http.MethodGet, url, nil, resp, opts...)
}

// Post http request
func Post(ctx context.Context, url string, req, resp interface{}, opts ...ReqOptionFunc) error {
	return HttpRequest(ctx, http.MethodPost, url, req, resp, opts...)
}

// Put http request
func Put(ctx context.Context, url string, req, resp interface{}, opts ...ReqOptionFunc) error {
	return HttpRequest(ctx, http.MethodPut, url, req, resp, opts...)
}

// Patch http request
func Patch(ctx context.Context, url string, req, resp interface{}, opts ...ReqOptionFunc) error {
	return HttpRequest(ctx, http.MethodPatch, url, req, resp, opts...)
}

// Delete http request
func Delete(ctx context.Context, url string, resp interface{}, opts ...ReqOptionFunc) error {
	return HttpRequest(ctx, http.MethodDelete, url, nil, resp, opts...)
}

// gin style log
func requestLog(printfer PrintfFunc, statusCode int, cost time.Duration, method, url string) {
	printfer("REQUEST | %3d | %13v | %-7s %s\n", statusCode, cost, method, url)
}

// http request
// req, resp are point type
func HttpRequest(ctx context.Context, method, addr string, req, resp interface{}, opts ...ReqOptionFunc) error {
	u, err := url.Parse(addr)
	if err != nil {
		return fmt.Errorf("bad request url: s%s", err)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	var (
		opt     = defaultReqOption
		reqBody io.Reader
	)
	for _, optFunc := range opts {
		if err := optFunc(&opt); err != nil {
			return err
		}
	}

	if opt.logTimeCost {
		start := time.Now()
		defer func() { requestLog(opt.printfer, *opt.statusCode, time.Since(start), method, addr) }()
	}

	if req != nil {
		var reqData []byte
		switch v := req.(type) {
		case *[]byte:
			reqData = *v
		case *string:
			reqData = []byte(*v)
		default:
			if reqData, err = opt.marshaler(req); err != nil {
				return err
			}
		}
		reqBody = bytes.NewReader(reqData)
	}

	request, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return err
	}

	nctx, cancel := context.WithTimeout(ctx, opt.timeout)
	defer cancel()
	request = request.WithContext(nctx)

	for k, v := range opt.header {
		request.Header.Set(k, v)
	}

	if opt.dumpRequest {
		CurlLikeDumpRequest(request)
	}

	response, err := opt.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if opt.dumpResponse {
		CurlLikeDumpResponse(response)
	}

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	*opt.statusCode = response.StatusCode
	if !opt.codeJudger(*opt.statusCode) {
		return NewRequestError(*opt.statusCode, string(respData))
	}
	if resp != nil {
		switch v := resp.(type) {
		case *[]byte:
			*v = respData
		case *string:
			*v = string(respData)
		default:
			err = opt.unmarshaler(respData, v)
			if err != nil {
				return fmt.Errorf("unmarshal %s failed: %s", respData, err)
			}
		}
	}

	return nil
}
