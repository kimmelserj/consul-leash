# Consul leash

[![Build Status](https://travis-ci.org/kimmelserj/consul-leash.svg?branch=master)](https://travis-ci.org/kimmelserj/consul-leash)
[![Coverage Status](https://coveralls.io/repos/github/kimmelserj/consul-leash/badge.svg?branch=master)](https://coveralls.io/github/kimmelserj/consul-leash?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kimmelserj/consul-leash)](https://goreportcard.com/report/github.com/kimmelserj/consul-leash)

Main destination of `consul-leash` is watch on specified consul key (environment variable `LEASH_KEY_PATH` value) value and compare consul key value with specified value (environment variable `LEASH_KEY_VALUE` value). If consul key value equals to specified value, then our application starts. If consul key value does not equal to specified value, then our application does not start and consul-leash waits until consul key value will be equal to specified value. If consul key value equals to specified value, then our application starts. If consul key value changed during your application works and consul key's new value does not equal specified value, consul-leash sends SIGTERM signal to your application and wait while your application going down. If your application hangs up, `consul-leash` waits specified period (environment variable `LEASH_STOPPING_DURATION` value) and kill your application and self going down with non-zero exit code.

Usually consul-leash is used as entry point for your applications. `consul-leash`'s stdout, stderr and stdin delegates to your application.

## Usage

    LEASH_KEY_PATH=my/daemon/current_master LEASH_KEY_VALUE=worker-1 consul-leash my-daemon arg1 arg2 arg3

## Docker image usage

    docker run --rm -it -e CONSUL_HTTP_ADDR=consul-agent:8500 -e LEASH_KEY_PATH=my/daemon/current_master -e LEASH_KEY_VALUE=worker-1 your-image-based-on-consul-leash my-daemon arg1 arg2 arg3
