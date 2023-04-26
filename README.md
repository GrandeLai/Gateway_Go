# Gateway_Go
gin+gorm+redis开发的微服务网关

实现http,tcp,grpc的反向代理和负载均衡等功能



启动方法

导入doc下面的sql文件

```
go mod tidy

go run main.go --config=./conf/dev --endpoint=dashboard

go run main.go --config=./conf/dev --endpoint=server
```

