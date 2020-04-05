package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func ping(ctx context.Context, s string) (interface{}, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ping")
	defer span.Finish()

	subPath := "/ping"
	method := "GET"
	userAgent := "hello-mongo"
	httpClient := &http.Client{}
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", s)
	}
	parsedURL.Path = path.Join(parsedURL.Path, subPath)

	req, err := http.NewRequest(method, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	// context
	req = req.WithContext(ctx)

	// headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	// Do Request
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// Parse Body
	var data interface{}
	if err := decodeBody(res, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
