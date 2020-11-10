#!/bin/bash
for ((;;))
do
  kubectl logs --follow keti-preprocessor-0 -n dt
  sleep 1
done


