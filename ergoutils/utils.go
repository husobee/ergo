// Copyright 2015 byteslice - all rights reserved
// This source code is governed by "The MIT License" which is found in the LICENSE file

// Package ergoutils - common utils for ergo
package ergoutils

import (
	"net/http"

	"golang.org/x/net/context"
)

const (
	// ContextErgoKey - the key by which one accesses the ergo struct from
	ContextErgoKey int = iota
	// ContextRequestKey - The key by which one accesses the http.Request from
	// the context.  request := ctx.Value(ContextRequestKey)
	ContextRequestKey
	// ContextMiddlewareStepCountKey - The key by which one accesses the Context Middleware
	ContextMiddlewareStepCountKey
)

// GetRequest - Helper function to get the request from the ctx
func GetRequest(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(ContextRequestKey).(*http.Request); ok && r != nil {
		return r
	}
	return nil
}
