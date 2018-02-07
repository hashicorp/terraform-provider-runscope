#!/usr/bin/env bash
set -ex
export VERSION=$(git describe --abbrev=0)
echo $VERSION
envsubst '$VERSION' < ./goreleaser.yml.template > ./goreleaser.yml
goreleaser --skip-validate