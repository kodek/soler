#!/bin/bash
git pull
./build_server.sh && \
  sudo systemctl restart soler-server

