package ergo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/byteslice/ergo/ergoutils"

	"log"
)

func TestMiddleware(t *testing.T) {
	// new ergo server
	e := NewErgo(context.Background())
	// Use the anonymous function as a middleware, which wraps sub middlewares
	e.Use(
		func(ctx context.Context, w http.ResponseWriter) error {
			// get current time
			start := time.Now()
			// log a message at start to see it wraps the sub "Next" handlerfunc
			log.Println("Start Middleware 1")

			// call the next function in the middleware chain
			err := Next(ctx, w)
			// check that the next, or nested middlewares, returned with an error
			if err != nil {
				log.Println("error from callee")
			}
			// print that the middleware finished, and how long it took
			log.Printf("End Middleware 1, took %d nanoseconds", time.Since(start).Nanoseconds())
			return nil

		},
		// Use the anonymous function as a middleware, which wraps sub middlewares
		func(ctx context.Context, w http.ResponseWriter) error {
			// get current time
			start := time.Now()
			// log a message at start to see it wraps the sub "Next" handlerfunc
			log.Println("Start Middleware 2")

			// call the next function in the middleware chain
			err := Next(ctx, w)
			// check that the next, or nested middlewares, returned with an error
			if err != nil {
				log.Println("error from callee")
			}
			// print that the middleware finished, and how long it took
			log.Printf("End Middleware 2, took %d nanoseconds", time.Since(start).Nanoseconds())
			return nil

		},
		// Use the anonymous function as a middleware, which wraps sub middlewares
		func(ctx context.Context, w http.ResponseWriter) error {
			// get current time
			start := time.Now()
			// log a message at start to see it wraps the sub "Next" handlerfunc
			log.Println("Start Middleware 3")

			// call the next function in the middleware chain
			err := Next(ctx, w)
			// check that the next, or nested middlewares, returned with an error
			if err != nil {
				log.Println("error from callee")
			}
			// print that the middleware finished, and how long it took
			log.Printf("End Middleware 3, took %d nanoseconds", time.Since(start).Nanoseconds())
			return nil

		})

	// Use this handler function, (same signature as a middlware)
	e.Use(
		func(ctx context.Context, w http.ResponseWriter) error {
			// print we are in the handler
			log.Println("in handler")
			// print we got the correct uri path from the request
			log.Printf("request uri: %s", ergoutils.GetRequest(ctx).URL.Path)
			return nil
		})

	// create a new request
	r, _ := http.NewRequest("GET", "/test", nil)
	// new response writer
	w := httptest.NewRecorder()
	// run the ergo server, with request and response recorder
	e.ServeHTTP(w, r)

	// create a new request
	r, _ = http.NewRequest("GET", "/test", nil)
	// new response writer
	w = httptest.NewRecorder()
	// run the ergo server, with request and response recorder
	e.ServeHTTP(w, r)
}
