package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.elastic.co/apm/module/apmhttp"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer("client-example")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

var header string

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}

	backend, ok := extra["backend"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}

	header, ok = extra["header"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return apmhttp.Wrap(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		req.Header.Set("Test",header)

		//Initiate the backend call via proxy
		serveReverseProxy(backend + req.URL.Path, w, apmhttp.RequestWithContext(ctx, req))
		
	})), nil
}

func init() {
	fmt.Println("client-example plugin loaded!!!")
}

func main() {}

//
// Reverse Proxy
//

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	//https://github.com/golang/go/issues/28168
	proxy := &httputil.ReverseProxy{
		//Transport: roundTripper(rt),
		Director: func(req *http.Request){
			req.Header.Set("Inline-Test", header)
			// Update the headers to allow for SSL redirection
			req.URL.Host = url.Host
			req.URL.Scheme = url.Scheme
			req.Host = url.Host
		},
	}

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}


//Custom roundtripper to be used in the reverse proxy.
func rt(req *http.Request) (*http.Response, error) {
    fmt.Printf("request received. url=%s", req.URL)
    req.Header.Set("Inline-Test", "headerwithinproxy")
    defer fmt.Printf("request complete. url=%s", req.URL)

    return http.DefaultTransport.RoundTrip(req)
}


// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
