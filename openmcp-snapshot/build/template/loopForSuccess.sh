#!/bin/bash

while :
do
  if [ -f "/success" ] ; then
    echo "success!!!! Job Exit"
    exit 0
  fi
  sleep 5
done