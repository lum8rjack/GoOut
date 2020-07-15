# Build using the following command
# docker build -t goout:1.0 .

# Run and compile the binary using the following command
# docker run --rm -v $(pwd)/compiled:/app/compiled/ goout:1.0

FROM golang:1.14
RUN mkdir /app
ADD . /app/
WORKDIR /app

# Build for Linux
CMD env GOOS=linux GOARCH=amd64 go build -o compiled/GoOut-Linux64 && chmod 777 compiled/*

# Build for Windows
#CMD env GOOS=windows GOARCH=amd64 go build -o build/GoOut-Windows64.exe && chmod 777 compiled/*

# Build for Mac
#CMD env GOOS=darwin GOARCH=amd64 go build -o build/GoOut-Darwin64 && chmod 777 compiled/*

# Build for Arm
#CMD env GOOS=linux GOARCH=arm go build -o build/GoOut-LinuxArm && chmod 777 compiled/*
