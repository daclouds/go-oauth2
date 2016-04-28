# Go Study

- Twelve Go Best Practices | Francesc Campoy Flores, Gopher at Google
    - https://talks.golang.org/2013/bestpractices.slide
    - https://gist.github.com/pzurek/6642797
- Best practices for a new Go developer | Satish Talim
    - https://blog.rubylearning.com/best-practices-for-a-new-go-developer-8660384302fc

### Getting started

##### redis (docker)

- https://hub.docker.com/_/redis/

```
$ docker run -d redis -p 6379:6379
```

##### Run

```
$ go run OAuth2.go
```

##### Test

Open in your web browser:
```
http://localhost:14000/authorize?response_type=code&client_id=1234&redirect_uri=http://localhost:14000/appauth/code
```

### OAuth 2 server library

- https://github.com/RangelReale/osin

### Redis client

- http://gopkg.in/redis.v3

### http

- https://golang.org/pkg/net/http
