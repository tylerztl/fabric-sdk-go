## Hyperledger Fabric Client SDK for Go
This GRPC server enables Go developers to build solutions that interact with Hyperledger Fabric by fabric-sdk-go .

### Prerequisites
- Go 1.10+ installation or later
- GOPATH environment variable is set correctly
- Govendor version 1.0.9 or later
- Protoc Plugin
- Protocol Buffers

### Getting started
Download fabric images
```
./scripts/download_images.sh
```
Start the fabric network
```
make networkUp
```
Start the fabric-sdk-go server
```
make start
```
Running the test suite
```
cd test
```

### Directory description
```
|--artifacts            // 启动fabric网络相关配置文件
|--conf                 // 项目配置文件
|  |--src               // 链码所在路径
|--logs                 // 日志文件所在路径
|--protos               // grpc所用的协议定义路径
|--scripts              // 脚本文件路径
|--server               // 服务功能实现
|  |--grpchandler       // 处理客户端grpc请求响应
|  |--helpers           // 通用的代码块
|  |--sdkprovider       // 调用fabric-sdk-go的接口
|  |--test              // 单元测试
|  |--vendor            // 项目依赖包 
```

