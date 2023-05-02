package httputil

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Give struct {
	method  string
	req     interface{}
	resp    interface{}
	timeout time.Duration
}

// go test -v -run TestHttpRequest
func TestHttpRequest(t *testing.T) {
	commonStr := `{"foo": "bar"}`
	commonData := []byte(commonStr)
	commonReq := make(map[string]string)
	json.Unmarshal(commonData, &commonReq) //nolint:errcheck
	b3Header := make(http.Header)
	b3Header.Set("x-b3-traceid", "abc")

	handler := func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(reqBody) //nolint:errcheck
			return
		}
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write(commonData) //nolint:errcheck
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write(reqBody) //nolint:errcheck
		case http.MethodPut:
			w.WriteHeader(http.StatusInternalServerError)
		case http.MethodPatch:
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			w.Write(commonData) //nolint:errcheck
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	tests := []struct {
		give     *Give
		wantCode int
		wantErr  bool
	}{
		{
			give: &Give{
				method: http.MethodGet,
				resp:   &[]byte{},
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			give: &Give{
				method: http.MethodPost,
				req:    &commonReq,
				resp:   &map[string]string{},
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			give: &Give{
				req:    &commonData,
				method: http.MethodPut,
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
		{
			give: &Give{
				method:  http.MethodPatch,
				timeout: 100 * time.Millisecond,
			},
			wantCode: http.StatusOK,
			wantErr:  true,
		},
		{
			give: &Give{
				method: http.MethodDelete,
				req:    &commonStr,
				resp:   new(string),
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give.method, func(t *testing.T) {
			d := 30 * time.Second
			if tt.give.timeout > 0 {
				d = tt.give.timeout
			}
			var statusCode int
			err := HttpRequest(context.TODO(), tt.give.method, ts.URL,
				tt.give.req, tt.give.resp,
				WithLogTimeCost(t.Logf),
				WithTimeout(d),
				WithClient(&http.Client{}),
				WithMarshal(json.Marshal),
				WithUnmarshal(json.Unmarshal),
				WithHeader(map[string]string{
					"Accept":       "application/json",
					"Content-Type": "application/json;charset=UTF-8",
				}),
				WithTraceHeaders(b3Header),
				WithStatusCodeJudge(defaultCodeJudger),
				StoreStatusCode(&statusCode))
			if tt.wantErr {
				t.Logf("err: %+v\n", err)
				code, ok := GetErrorCode(err)
				if ok {
					assert.Equal(t, code, tt.wantCode)
					// for test cover
					if !IsErrorCode(err, tt.wantCode) {
						t.Fatalf("%+v wanted status code: %d", tt.give, tt.wantCode)
					}
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, statusCode, tt.wantCode)
				switch v := tt.give.resp.(type) {
				case *map[string]string:
					assert.Equal(t, (*v)["foo"], "bar")
				case *[]byte:
					assert.Equal(t, *v, commonData)
				case *string:
					assert.Equal(t, *v, commonStr)
				}
			}
		})
	}
}
