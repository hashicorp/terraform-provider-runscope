#!/usr/bin/env bash
set -ex
VERSION=$(git describe --abbrev=0)
find dist/ -name terraform-provider-runscope_$VERSION | xargs chmod +x