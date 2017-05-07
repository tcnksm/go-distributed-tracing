# Message service protobuf

This repository contains protobuf message service definition.

To generate go package,

```bash
$ protoc -I=./ --go_out=plugins=grpc:./ message.proto
```
