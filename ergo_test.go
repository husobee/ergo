package ergo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log"

	"golang.org/x/net/context"
)

func TestMiddleware(t *testing.T) {
	// new ergo server
	e := NewErgo()
	// Use the anonymous function as a middleware, which wraps sub middlewares
	e.Use(
		func(ctx context.Context, w http.ResponseWriter) error {
			// get current time
			start := time.Now()
			// log a message at start to see it wraps the sub "Next" handlerfunc
			log.Println("Start Middleware")

			// call the next function in the middleware chain
			err := Next(ctx, w)
			// check that the next, or nested middlewares, returned with an error
			if err != nil {
				log.Println("error from callee")
			}
			// print that the middleware finished, and how long it took
			log.Printf("End Middleware, took %d nanoseconds", time.Since(start).Nanoseconds())
			return nil

		})

	// Use this handler function, (same signature as a middlware)
	e.Use(
		func(ctx context.Context, w http.ResponseWriter) error {
			// print we are in the handler
			log.Println("in handler")
			// print we got the correct uri path from the request
			log.Printf("request uri: %s", GetRequest(ctx).URL.Path)
			return nil
		})

	// create a new request
	r, _ := http.NewRequest("GET", "/test", nil)
	// new response writer
	w := httptest.NewRecorder()
	// run the ergo server, with request and response recorder
	e.ServeHTTP(w, r)
}
