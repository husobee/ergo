// Copyright 2015 byteslice - all rights reserved
// This source code is governed by "The MIT License" which is found in the LICENSE file

// Package ergo - ergo is a simple and sane web application framework
package ergo

import (
	"errors"
	"net/http"

	"github.com/byteslice/ergo/ergoutils"

	"golang.org/x/net/context"
)

var (
	// ErrNoMiddleware - Error when calling Next to get next middleware handler
	ErrNoMiddleware = errors.New("Ergo - No Next Middleware")
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

//Middleware - Interface to define what a Middleware is
type Middleware interface {
	GetFunc(int) (HandlerFunc, error)
}

// middleware is a structure to hold chained middlewares and positional indication
type middleware struct {
	chain []HandlerFunc
}

// middleware counter struct
type counter struct {
	i int
}

// GetFunc - implementation of Middleware
func (m *middleware) GetFunc(i int) (HandlerFunc, error) {
	if len(m.chain) >= i {
		return m.chain[i-1], nil
	}
	return nil, ErrNoMiddleware
}

// Use - Insert a middleware into the chain
func (e *Ergo) Use(h HandlerFunc) {
	if m, ok := e.context.Value(ergoutils.ContextMiddlewareKey).(*middleware); ok && m != nil {
		m.chain = append(m.chain, h)
		return
	}
	e.context = context.WithValue(e.context, ergoutils.ContextMiddlewareKey, &middleware{chain: []HandlerFunc{h}})
}

// ServeHTTP - implementation of standard http.Handler
func (e *Ergo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// define new context from base context, include request
	ctx := context.WithValue(e.context, ergoutils.ContextRequestKey, r)
	ctx = context.WithValue(ctx, ergoutils.ContextMiddlewareCounterKey, &counter{0})
	Next(ctx, w)
}

// Next - Helper function to call the next middleware from the current handlerfunc
func Next(ctx context.Context, w http.ResponseWriter) error {
	if m, ok := ctx.Value(ergoutils.ContextMiddlewareKey).(Middleware); ok && m != nil {
		// increment counter off context
		if counter, ok := ctx.Value(ergoutils.ContextMiddlewareCounterKey).(*counter); ok {
			counter.i++
			f, err := m.GetFunc(counter.i)
			if err != nil {
				if err == ErrNoMiddleware {
					// end of chain
					return nil
				}
				return err
			}
			if err := f(ctx, w); err != nil {
				return err
			}

			return Next(ctx, w)
		}
	}
	return ErrNoMiddleware
}
