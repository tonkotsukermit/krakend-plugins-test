FROM golang:1.14-alpine as build
WORKDIR /go/src/app
COPY . .

#CGO disabled gives a static binary, OS set to linux, -a rebuild all packages
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app ./cmd/echo

FROM scratch
# the app
COPY --from=build /go/bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]