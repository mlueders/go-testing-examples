#!/bin/bash

for branch in $(git for-each-ref --format='%(refname)' refs/heads/); do
    branchName=`echo "${branch##*/}"`
    if [ "${branchName}" != "main" ]; then
      git checkout "${branchName}"
      git rebase main
      if [ $? != 0 ]; then
        echo "rebase failed"
        exit
      fi
    fi
done

git checkout main
