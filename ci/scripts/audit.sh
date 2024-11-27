#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/s3crypto
  make audit
popd