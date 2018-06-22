# Consul leash

## Usage

    LEASH_KEY_PATH=my/daemon/current_master LEASH_KEY_VALUE=worker-1 consul-leash my-daemon arg1 arg2 arg3

## Docker image usage

    docker run --rm -it -e CONSUL_HTTP_ADDR=consul-agent:8500 -e LEASH_KEY_PATH=my/daemon/current_master -e LEASH_KEY_VALUE=worker-1 your-image-based-on-consul-leash my-daemon arg1 arg2 arg3
