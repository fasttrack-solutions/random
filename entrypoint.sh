#!/bin/sh
set -e

if [ "$1" = "grpc" ]; then
  exec ./grpc
elif [ "$1" = "http" ]; then
  exec ./http
else
  echo "Usage: docker run <image> [grpc|http]"
  exit 1
fi
