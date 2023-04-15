#!/bin/bash

gcc -g main.c -o main -lavformat -lavcodec -lavutil
./main
