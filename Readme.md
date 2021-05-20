# KrakenD Plugin Test

This is a simple project to test krakend plugins

## Compose

Docker compose is used to build the plugins, start up krakend and the echo service. `docker-compose up`

If you run into issues with plugins not reloading properly in the compose, be sure to run `docker system prune -af`  to clear all docker artifacts if needed, and `docker-compose build` to rebuild the images

## Building Plugins

The plugins are currently built within the `./build/krakend/Dockerfile`. This is to ensure the dependancies of the plugins as well as the build of kraken don't have any conflicts. Plugins are currently not supported in windows, so if you wish to build them on windows, you will need to run them in WSL.

linux:
`go build -buildmode=plugin -o ./plugin/header-example.so ./router/header`

### Client (proxy)

intercepts the call before kraken can do anything with it.

```go
func (r registerer) RegisterClients(f func(
    name string,
    handler func(context.Context, map[string]interface{}) (http.Handler, error),
))
```

### Server (router)

Changes the call midflight

```go
func (r registrable) RegisterHandlers(f func(
name string,
handler func(
context.Context,
map[string]interface{},
http.Handler) (http.Handler, error),
)) {
f(string(r), r.registerHandlers)
}
```

## Usage

This config resides at the highest level of extra config before the endpoints are called. [ref](https://www.eventslooped.com/posts/krakend-writing-plugins/)

Example:

```json
{
    "version": 2,
    "extra_config": {
      "github_com/devopsfaith/krakend/transport/http/server/handler": {
        "name": "header-example",
        "extraconfig":"something"
     }

    },
    "timeout": "3000ms",
    "cache_ttl": "300s",
    "output_encoding": "json",
    "name": "test-service",
    "plugin": {
        "pattern":".so",
        "folder": "./plugin/"
    },
    "endpoints": []
    }
```

there is a way to implement a client executor within the endpoint's extra config as shown [here](https://www.krakend.io/blog/krakend-grpc-gateway-plugin/) by implementing your own http client.
