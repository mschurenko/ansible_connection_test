#!/bin/bash

set -e

for i in http port;do
  if [ $(library/check_${i}_osx check_${i}/test.json|jq -r '.failed') != false ];then
    echo check_${i} failed
    exit 1
  fi
done

