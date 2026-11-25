package space

import (
	"context"
	"io"
	"net/http"

	"github.com/reiver/go-xrpc"
)

// ContentTypeSTAR is the media type of a STAR file (ex: from getRepo, exportSpaceMeta; to importRepo, importSpaceMeta).
const ContentTypeSTAR = "application/vnd.dasl.star"

func query[T any](ctx context.Context, client xrpc.Client, url string, header http.Header) (T, error) {
	var output T
	if err := client.QueryJSON(ctx, &output, url, header); nil != err {
		var empty T
		return empty, err
	}
	return output, nil
}

func execute[T any](ctx context.Context, client xrpc.Client, url string, header http.Header, input any) (T, error) {
	var output T
	if err := client.ExecuteJSON(ctx, &output, url, header, input); nil != err {
		var empty T
		return empty, err
	}
	return output, nil
}

// executeNoOutput calls a procedure that has no output (or an empty one).
func executeNoOutput(ctx context.Context, client xrpc.Client, url string, header http.Header, input any) error {
	return client.ExecuteJSON(ctx, nil, url, header, input)
}

// executeSTAR calls a procedure whose input is a STAR file, and which has no output.
func executeSTAR(ctx context.Context, client xrpc.Client, url string, header http.Header, body io.Reader) error {
	response, err := client.Execute(ctx, url, header, ContentTypeSTAR, body)
	if nil != err {
		return err
	}
	return response.Body.Close()
}
