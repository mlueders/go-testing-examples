#!/bin/bash

mkdir reports
rm reports/*

export SHOULD_FAIL=true
for branch in $(git for-each-ref --format='%(refname)' refs/heads/); do
    branchName=`echo "${branch##*/}"`
    git checkout $branch
    go test -v > "reports/${branchName}.out"
done

git checkout main
