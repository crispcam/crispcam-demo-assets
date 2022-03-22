package crisps

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type contextKey string

var (
	contextKeyHeaders = contextKey("headers")
)

func (c contextKey) String() string {
	return "crisps context key: " + string(c)
}

// TraceRequest Persist Istio Tracing Headers
func TraceRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headers = http.Header{}
		ctx := r.Context()
		tracingHeaders := []string{
			"x-request-id",
			"x-b3-traceid",
			"x-b3-spanid",
			"x-b3-sampled",
			"x-b3-parentspanid",
			"x-b3-flags",
			"x-ot-span-context",
			"user-agent",
		}
		for _, key := range tracingHeaders {
			if val := r.Header.Get(key); val != "" {
				// Persist headers for both GRPC and HTTP
				ctx = metadata.AppendToOutgoingContext(ctx, key, val)
				headers.Add(key, val)
			}
		}
		ctx = context.WithValue(ctx, contextKeyHeaders, headers)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TraceHeaders(ctx context.Context) (http.Header, bool) {
	headers, ok := ctx.Value(contextKeyHeaders).(http.Header)
	return headers, ok
}

func Request(r *http.Request, u string, method string, form url.Values) ([]byte, error) {
	var result []byte
	req, err := http.NewRequest(method, u, strings.NewReader(form.Encode()))
	if err != nil {
		return result, err
	}
	// Persist http headers
	req = req.WithContext(r.Context())
	headers, ok := TraceHeaders(r.Context())
	// Don't worry if this doesn't work, just don't persist them
	if ok {
		req.Header = headers
	}
	// Assume json if there is no body, otherwise it's an encoded form
	if form == nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	} else {
		req.PostForm = form
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("upstream status code %d (request URI: %v)", resp.StatusCode, u))
	}
	result, err = ioutil.ReadAll(resp.Body)
	return result, nil
}
