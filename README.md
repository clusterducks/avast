# avast
...ye maties!

![Avast](https://media.giphy.com/media/E8KFBhPh2s3ra/giphy.gif)

## Server-side

### Building

```
get get
# recommended:
# gb vendor fetch --all

go build -o bin/avast github.com/bfowle/avast/src/
```

^-- `@TODO` change to **make**

### Running

```
AVAST_API_VERSION=v1 \
  AVAST_ADDR=:8080 \
  DOCKER_HOST=tcp://123.45.67.890:2375 \
  DOCKER_API_VERSION=v1.21 \
  CONSUL_HTTP_ADDR=123.45.67.890:8500 \
  bin/avast
```
  
## Client-side

### Building

```
cd client && npm i
```

^-- `@TODO` change to **make**

### Running

```
npm start # in client/
```

## Production

TBD--
