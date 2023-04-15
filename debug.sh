#!/bin/bash

if ! gcc -g main.c -o main -lavformat -lavcodec -lavutil; then
    echo "Compilation failed"
    exit 1
fi

if ! valgrind --leak-check=full --show-leak-kinds=all ./main; then
    echo "Memory leak detected"
    exit 1
fi
