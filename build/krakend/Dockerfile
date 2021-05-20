# Note: If you encounter any issues with package mismatch errors
# Ensure that the version of Go used to compile is the same version as Go for KrakenD.
# Search in the github.com/devopsfaith/krakend-ce repo for GOLANG_VERSION to check this.
FROM golang:1.14.0 AS go
WORKDIR /opt/
RUN git clone https://github.com/devopsfaith/krakend-ce.git .
RUN git checkout tags/v1.1.1
COPY ./plugins ./plugins

#Ensure that krakends dependancies are the same as the plugins to avoid compatibility issues
RUN go mod tidy

#Build the plugins here
RUN go build -buildmode=plugin -o /opt/out/client-example.so /opt/plugins/client-example
RUN go build -buildmode=plugin -o /opt/out/handler-example.so /opt/plugins/handler-example

RUN make build

FROM devopsfaith/krakend:1.1.1
WORKDIR /etc/krakend
#Copy the plugins from the build to the kraken image
COPY --from=go /opt/out/*.so /etc/krakend/plugins/
COPY --from=go /opt/krakend /usr/bin/krakend
