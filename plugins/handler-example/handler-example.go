package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.elastic.co/apm/module/apmhttp"
)

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer("handler-example")

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return apmhttp.Wrap(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		//Create a new request to pass back to kraken with the extra header values injected
		r2 := new(http.Request)

		//Populate the new request with all previous request values.
		*r2 = *req

		r2.Header.Add("some","test")

		fmt.Printf(req.RequestURI)

		
		handler.ServeHTTP(w, r2)
		
	})), nil
}

func init() {
	fmt.Println("handler-example handler plugin loaded!!!")
}

func main() {}