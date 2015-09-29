// Copyright 2015 byteslice - all rights reserved
// This source code is governed by "The MIT License" which is found in the LICENSE file

// package ergo - ergo is a simple and sane web application framework
package ergo

import (
	"errors"
	"net/http"

	"golang.org/x/net/context"
)

const (
	// ContextRequestKey - The key by which one accesses the http.Request from
	// the context.  request := ctx.Value(ContextRequestKey)
	contextRequestKey int = iota
	middlewareKey
	middlewareCounter
)

// HandlerFunc - Ergo Handler Function type
type HandlerFunc func(context.Context, http.ResponseWriter) error

// Ergo - Main Ergo structure, base context by which everything is derived
type Ergo struct {
	context context.Context
}

// NewErgo - Create a new *Ergo structure
func NewErgo() *Ergo {
	return &Ergo{
		context: context.Background(),
	}
}

// middleware is a structure to hold chained middlewares and positional indication
type middleware struct {
	chain   []HandlerFunc
	counter int
}

// GetRequest - Helper function to get the request from the ctx
func GetRequest(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(contextRequestKey).(*http.Request); ok && r != nil {
		return r
	}
	return nil
}

// Next - Helper function to call the next middleware from the current handlerfunc
func Next(ctx context.Context, w http.ResponseWriter) error {
	if m, ok := ctx.Value(middlewareKey).(*middleware); ok && m != nil {
		m.counter += 1
		return m.chain[m.counter-1](ctx, w)
	}
	return errors.New("middleware chain break")
}

// Use - Insert a middleware into the chain
func (e *Ergo) Use(h HandlerFunc) {
	if m, ok := e.context.Value(middlewareKey).(*middleware); ok && m != nil {
		m.chain = append(m.chain, h)
		return
	}
	e.context = context.WithValue(e.context, middlewareKey, &middleware{chain: []HandlerFunc{h}, counter: 0})
}

// ServeHTTP - implementation of standard http.Handler
func (e *Ergo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// define new context from base context, include request
	ctx := context.WithValue(e.context, contextRequestKey, r)
	Next(ctx, w)
}
