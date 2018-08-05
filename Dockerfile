FROM alpine:3.7

COPY consul-leash /consul-leash

ENTRYPOINT ["/consul-leash"]
