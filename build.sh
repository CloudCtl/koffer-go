#!/bin/bash -x
sudo /usr/bin/podman run \
    -it --rm --name go-build \
    --volume $(pwd)/bin:/tmp/bin:z \
    --entrypoint /root/dev/dev.sh \
    --volume $(pwd):/root/dev:z \
  docker.io/ocpredshift/red-gotools
    
