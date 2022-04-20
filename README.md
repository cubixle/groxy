# Groxy

Groxy is a lightweight, simple HTTP reverse proxy written in golang.

Groxy was built to run on port 8080 by default. This can be changed with the the `GROXY_PORT` environment variable.

## Config

The config for groxy is very simple and will hopefully always be simple.

```yaml
debug: false
endpoints: 
  - addr: pihole.test
    remote_addr: "192.168.1.2:8017"
  - addr: trans.test
    remote_addr: "192.168.1.2:9091"
```

By Default groxy will try and read a `groxy.yml` file from its current working directory. To override this you can use the environment variable `GROXY_CONFIG_FILE`.

## Metrics

There are a bunch of different metrics that can be found on port 9090.

```
# TYPE request_counter counter
request_counter{code="200",method="GET"} 209
request_counter{code="304",method="GET"} 20
```

## Docker

Docker containers can be found [https://hub.docker.com/r/cubixle/groxy](here)

