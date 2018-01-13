# Usage

1. Start etcd

```
./etcd
```

2. Build tools

```
go build -o server ../cmd/server.go
go build -o agent ../cmd/agent.go
go build -o soberctl ../cmd/soberctl.go
```

3. Initial authorize

```
./soberctl set /prod/infrastructure/service/sober
./soberctl get /prod/infrastructure/service/sober
```

4. Initial project

```
./soberctl set /prod/blog/backend/go
./soberctl get /prod/blog/backend/go
```

5. Start server

```
./server -debug
```

6. Start agent

```
./agent --datasource grpc://127.0.0.1:3333 --root /prod/blog/backend/go --token prodgotoken --output file://prodBlogBackendGo.json --debug
```    