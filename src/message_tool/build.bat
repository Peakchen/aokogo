cd protoc
protoc -I=../message/ --go_out=../go/ c2s_message.proto
protoc -I=../message/ --go_out=../go/ error_message.proto
protoc -I=../message/ --go_out=../go/ s2s_message.proto

pause