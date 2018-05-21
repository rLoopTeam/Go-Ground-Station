protoc --go_out=plugins=grpc:. groundstation.proto
python -m grpc_tools.protoc -I. --python_out=./python_proto --grpc_python_out=./python_proto groundstation.proto