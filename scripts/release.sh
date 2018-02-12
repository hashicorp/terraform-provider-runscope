#!/usr/bin/env bash
set -ex
export VERSION=$(git describe --abbrev=0)
echo $VERSION
goreleaser --config <(envsubst '$VERSION' < ./goreleaser.yml.template)