#!/usr/bin/env bash

docker run -d \
  --name ha \
  --restart=unless-stopped \
  -e TZ=Europe/Brussels \
  -v ${PWD}/config:/config \
  -p 8123:8123 \
  ghcr.io/home-assistant/home-assistant:stable

  #  --network=host \
  #--privileged \