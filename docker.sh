#!/bin/bash
docker run \
    -it --rm --net=host \
    --volume $(pwd):/root/dev \
    --volume ~/.ssh:/root/.ssh \
    --volume ~/.bashrc:/root/.bashrc \
    --volume ~/.gitconfig:/root/.gitconfig \
    --name koffer-go --hostname koffer-go \
  docker.io/ocpredshift/red-gotools -c /usr/bin/tmux
