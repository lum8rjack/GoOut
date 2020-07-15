#!/bin/bash

# Build using the following command
docker build -t goout-server:1.0 .

# Run and compile the binary using the following command
docker run --rm -v $(pwd)/compiled:/app/compiled/ goout-server:1.0

# Copy the binary to the cwd since it needs access to the other locations
cp compiled/* .