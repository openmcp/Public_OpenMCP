#!/bin/bash

pyinstaller -F ketikubecli.py
cp dist/ketikubecli /bin
