# Go-Ground-Station

The groundstation's backend written in GO.

This part of the groundstation is the side that communicates directly with the rPod/rFlight etc...

The groundstation uses the gRPC protocol to communicate with whatever type of frontend we'd choose to use.

# Set up

Before pulling the repo make sure you understand [how Go development works](https://golang.org/doc/code.html#Organization)

The project is structured as follows: $GOPATH/src/rloop/<repo>

# Dependencies

#### gRPC

According to the [gRPC website](https://grpc.io/docs/quickstart/go.html):

- Make sure you are using Go 1.6+ ```$ go version```
- Install gRPC  ```$ go get -u google.golang.org/grpc```
- Install Protocol Buffers v3:
  - First you need to [download the pre-compiled binaries for the protoc compiler](https://github.com/google/protobuf/releases), 
  unzip the files and add the path to the folder in your environment variables
  - Next install the protoc plugin for Go ```$ go get -u github.com/golang/protobuf/protoc-gen-go```


# Running Locally

While developing you can simply run the app with the next command: ```$ go run main.go```

# Compile

# Available Scripts

###### vet_all
Vet examines Go source code and reports suspicious constructs, [...] it can find errors not caught by the compilers.

###### build
The build script is useful as it will compile the go code for the 3 main platforms (Windows,Mac and Linux)
