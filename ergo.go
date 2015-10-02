// Copyright 2015 byteslice - all rights reserved
// This source code is governed by "The MIT License" which is found in the LICENSE file

// Package ergo - ergo is a simple and sane web application framework
package ergo

import (
	"errors"
	"net/http"
	"sync"

	"github.com/byteslice/ergo/ergoutils"

	"golang.org/x/net/context"
)

var (
	// ErrNoMiddleware - Error when calling Next to get next middleware handler
	ErrNoMiddleware       = errors.New("Ergo - No Next Middleware")
	ErrEndMiddlewareChain = errors.New("Ergo - End of Middleware Chain")
)

// HandlerFunc - Ergo Handler Function type, this is the base structure for handlers
// and middlewares within ergo, defines a context as a parameter, and an
// http.ResponseWriter as a parameter.  Error responses will be picked up in the
// error handling middleware
type HandlerFunc func(context.Context, http.ResponseWriter) error

// Ergo - Main Ergo structure, base context and middleware chain.  The middleware
// chain will consist of a MiddlewareChain implementation
type Ergo struct {
	context    context.Context
	middleware MiddlewareChain
}

// NewErgo - Create a new *Ergo structure
func NewErgo(c context.Context) *Ergo {
	// create ergo pointer with empty middleware
	e := &Ergo{
		middleware: newMiddleware(),
	}
	if c == nil {
		// if user doesn't pass in a context to ergo, use base background context
		e.context = context.WithValue(context.Background(), ergoutils.ContextErgoKey, e)
	} else {
		// if user passes in a context to ergo, use existing context
		e.context = context.WithValue(c, ergoutils.ContextErgoKey, e)
	}
	return e
}

//Middleware - Interface to define what a Middleware is
type MiddlewareChain interface {
	// GetFunc - This function will get the handler at the position defined by
	// the int parameter, and respond with a handler function and error.
	GetFunc(int) (HandlerFunc, error)
	// AddFunc - This function will allow for additions of middlewares to ergo
	AddFunc(HandlerFunc)
}

// middleware is a structure to hold chained middlewares and positional indication
type middleware struct {
	m     *sync.RWMutex
	chain []HandlerFunc
}

// newMiddleware - create a new middleware structure
func newMiddleware() *middleware {
	return &middleware{
		m:     new(sync.RWMutex),
		chain: []HandlerFunc{},
	}
}

// middleware step counter struct
type stepCount struct {
	i int
}

// GetFunc - implementation of Middleware
func (m *middleware) GetFunc(i int) (HandlerFunc, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	if len(m.chain) >= i {
		return m.chain[i-1], nil
	}
	return nil, ErrEndMiddlewareChain
}

// AddFunc - implementation of Middleware
func (m *middleware) AddFunc(h HandlerFunc) {
	m.m.Lock()
	defer m.m.Unlock()
	m.chain = append(m.chain, h)
}

// Use - Insert a middleware or many into the chain, processes in order
func (e *Ergo) Use(h ...HandlerFunc) {
	for _, v := range h {
		e.middleware.AddFunc(v)
	}
}

// ServeHTTP - implementation of standard http.Handler
func (e *Ergo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// define new context from base context, include request
	ctx := context.WithValue(e.context, ergoutils.ContextRequestKey, r)
	// define the context middlewareStepCountKey
	ctx = context.WithValue(ctx, ergoutils.ContextMiddlewareStepCountKey, &stepCount{0})
	// process the first function in the middlware chain, or handler
	// TODO: what if we get an error on this call?
	Next(ctx, w)
}

// Next - Helper function to call the next middleware from the current handlerfunc
func Next(ctx context.Context, w http.ResponseWriter) error {
	// from the context, grab the Ergo pointer, so we can access middleware chain
	if e, ok := ctx.Value(ergoutils.ContextErgoKey).(*Ergo); ok && e != nil {
		// pull middleware chain from ergo
		m := e.middleware
		// increment counter off context
		if counter, ok := ctx.Value(ergoutils.ContextMiddlewareStepCountKey).(*stepCount); ok {
			counter.i++
			// get the function
			f, err := m.GetFunc(counter.i)
			// if there is an error getting the middleware function
			if err != nil {
				if err == ErrEndMiddlewareChain {
					// end of chain
					return nil
				}
				return err
			}
			// perform the handler call, and return the error if recieved
			if err := f(ctx, w); err != nil {
				return err
			}
			// call the next middleware in the chain
			// usually this will be performed within the middlware
			// but this is a catch all in the event the middlewares do not
			// call Next
			return Next(ctx, w)
		}
	}
	return ErrNoMiddleware
}
