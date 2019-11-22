setlocal

@rem enter this directory
cd /d %~dp0

protoc --proto_path=protos --micro_out=src/helloworld --go_out=src/helloworld protos/greeter.proto

endlocal