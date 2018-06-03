# Usage

1. Start etcd

```
./etcd
```

2. Build tools

```
go build -o bin/server ../cmd/server/main.go
go build -o bin/agent ../cmd/agent/main.go
go build -o bin/soberctl ../cmd/soberctl/main.go
```

3. Initial authorize

```
bin/soberctl set /prod/infrastructure/service/sober.authorize.toml
bin/soberctl get /prod/infrastructure/service/sober.authorize.toml
```

4. Initial project

```
bin/soberctl set /prod/blog/backend/go.json
bin/soberctl get /prod/blog/backend/go.json
```

5. Start server

```
bin/server -debug
```

6. Start agent

```
bin/agent --datasource grpc://127.0.0.1:3333 --key /prod/blog/backend/go.json --token go --output file://prodBlogBackendGo.json --debug
```

