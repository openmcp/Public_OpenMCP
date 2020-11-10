#!/bin/bash
for (( ; ; ))
do
  kubectl logs --follow keti-ruleengine-0 -n dt
  sleep 1
done
