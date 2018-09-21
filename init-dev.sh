#!/bin/bash

# This script initializes the repo with the necessary Go deps and git hooks.
# Usage: `./init-dev.sh`

# Install all packages required by eventmanager.
# Please ensure you have `dep` installed locally first (see README.md)!
dep ensure

# Install all tools needed for eventmanager development.
go get -u github.com/golang/lint/golint
go get -u github.com/ashwch/precommit-vet-lint

# Install local git hooks for linting.
# Any new git hook scripts should be specified here.
ln -s -f ../../scripts/hooks/pre-push .git/hooks/pre-push
