# avast
...ye maties!

![Avast](https://media.giphy.com/media/E8KFBhPh2s3ra/giphy.gif)

### Dependencies

- `go` (1.4+)

### Building

```
go get
go build -o bin/avast
```

### Running

```
AVAST_API_VERSION=v1 \
  AVAST_ADDR=:8080 \
  AVAST_DATACENTER=dc1 \
  DOCKER_HOST=tcp://1.1.1.1:2378 \
  DOCKER_API_VERSION=v1.21 \
  CONSUL_HTTP_ADDR=1.1.1.1:8500 \
  bin/avast
```
