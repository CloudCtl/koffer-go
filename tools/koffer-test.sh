#!/bin/bash
# allows testing koffer runs against plugins with the pre-built binary
# Example:
# - $ ./tools/koffer-test.sh --repo collector-infra --branch master --silent

# Ensure artifact dir is present
mkdir -p /tmp/bundle ; \

# Run koffer
run_test () {
run_home=${HOME}
sudo podman run -it --rm --pull always \
    --volume ${run_home}/.docker:/root/.docker:z \
    --volume /tmp/bundle:/root/deploy/bundle:z \
    --volume $(pwd)/bin/koffer:/usr/bin/koffer:z \
  docker.io/containercraft/koffer bundle $@
}

run_test $@
