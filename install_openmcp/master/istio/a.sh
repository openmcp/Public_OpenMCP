#!/bin/bash

for num in `seq 1 10`
do 
    curl -v -HHost:webservice.greeting.local http://10.0.3.203:80/ | tail -1 >> ~/output.txt
done
