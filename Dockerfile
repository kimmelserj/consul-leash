FROM alpine:3.7

COPY cmd/consul-leash/consul-leash /consul-leash

ENTRYPOINT ["/consul-leash"]
