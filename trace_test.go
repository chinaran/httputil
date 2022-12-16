package httputil

import (
	"net/http"
	"reflect"
	"testing"
)

func TestForwardHeaders(t *testing.T) {
	type args struct {
		headers http.Header
	}

	b3Header := make(http.Header)
	b3Header.Set("x-b3-traceid", "abc")

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "b3",
			args: args{
				headers: b3Header,
			},
			want: map[string]string{
				"x-b3-traceid": "abc",
			},
		},
		{
			name: "empty",
			args: args{
				headers: map[string][]string{},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ForwardHeaders(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForwardHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
