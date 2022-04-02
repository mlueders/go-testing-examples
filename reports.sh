#!/bin/bash

rm -rf reports/
mkdir -p reports/pass
mkdir -p reports/fail

for branch in $(git for-each-ref --format='%(refname)' refs/heads/); do
    branchName=`echo "${branch##*/}"`
    git checkout $branch
    export SHOULD_FAIL=false
    go test -v > "reports/pass/${branchName}.out"
    export SHOULD_FAIL=true
    go test -v > "reports/fail/${branchName}.out"
done

git checkout main
