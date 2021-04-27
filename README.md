# genwasm
Generate wasm file by giving an input as go file

To run the wasm builder in docker for maintaining the go version. Enter the below command:

  docker run --rm -it --name wasmbuilder -v $PWD:/go/src/github.com/dhawalhost/genwasm -w /go/src/github.com/dhawalhost/genwasm golang:1.15.6 make all
