#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/s3crypto
  make test
popd