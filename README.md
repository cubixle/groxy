# Groxy

Groxy is a lightweight, simple HTTP reverse proxy written in golang.

Groxy was built to run in a docker container so is setup by default to listen on port 80.

## Config

The config for groxy is very simple and will hopefully always be simple.

```yaml
endpoints: 
  - addr: pihole.test
    remote_addr: "192.168.1.2:8017"
  - addr: trans.test
    remote_addr: "192.168.1.2:9091"
```

## Metrics

There are a bunch of different metrics that can be found on port 9090.

```
# TYPE request_counter counter
request_counter{code="200",method="GET"} 209
request_counter{code="304",method="GET"} 20
```

## Docker

Docker containers can be found here and are built for amd64 and arm64.

