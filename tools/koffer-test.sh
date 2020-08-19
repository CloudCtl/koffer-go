#!/bin/bash
# allows testing koffer runs against plugins with the pre-built binary
# Example:
# - $ ./tools/koffer-test.sh --repo collector-apps --branch master

# Ensure artifact dir is present
mkdir -p /tmp/bundle ; \

# Run koffer
run_test () {
sudo podman run -it --rm --pull always \
    --volume /tmp/bundle:/root/deploy/bundle:z \
    --volume $(pwd)/bin/koffer:/usr/bin/koffer:z \
  docker.io/codesparta/koffer bundle $@
}

run_test $@
