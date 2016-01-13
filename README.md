# avast
...ye maties!

![Avast](https://media.giphy.com/media/E8KFBhPh2s3ra/giphy.gif)

## Server-side

### Building

```go build -o bin/avast github.com/bfowle/avast/src/```

^-- TODO: change to make

### Running

```DOCKER_HOST=tcp://123.45.67.890:1234 \
  DOCKER_API_VERSION=v1.21 \
  CONSUL_HTTP_ADDR=123.45.67.890:8500 \
  bin/avast```
  
## Client-side

### Building

```npm install```

### Running

```npm start```

## Production

TBD--
