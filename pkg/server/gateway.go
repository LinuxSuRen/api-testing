package server

import (
	context "context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// MetadataStoreFunc is a function that extracts metadata from a request.
func MetadataStoreFunc(ctx context.Context, r *http.Request) (md metadata.MD) {
	store := r.Header.Get(HeaderKeyStoreName)
	md = metadata.Pairs(HeaderKeyStoreName, store)
	return
}
