#!/bin/bash
echo "[ketikubecli Copy] Path : /bin/ketikubecli"
cp dist/ketikubecli /bin


echo "[Config File Copy] Path : /var/lib/ketikubecli/config.yaml"
mkdir -p /var/lib/ketikubecli
cp config.yaml /var/lib/ketikubecli/config.yaml



