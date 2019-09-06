# prometheus-m3db-sarama
prometheus获取数据通过kafka存储到远端M3DB数据
## development

```
go test
go build
```

## docker run 
```
make build

docker build -t saramam3db:1.0.0 .

docker run -d --net=host --name saramam3db saramam3db:1.0.0  --store-url=http://10.254.192.2:7201/api/v1/prom/remote/write
```