#!/bin/bash

# Build using the following command
docker build -t goout:1.0 .

# Run and compile the binary using the following command
docker run --rm -v $(pwd)/compiled:/app/compiled/ goout:1.0