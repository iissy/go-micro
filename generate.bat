setlocal

@rem enter this directory
cd /d %~dp0

protoc -I../../../ -I=. --go_out=../../../ ./messages/messages.proto
protoc -I../../../ -I=. --micro_out=plugins=micro:. ./helloworld/helloworld.proto

endlocal