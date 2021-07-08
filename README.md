# Redis Cache package

This package provides data storage in the key value format in the radis database

To run the redisDB docker image, use
```shell
build/redis_up.sh
```

## Usage:

```go
import "github.com/rurick/rediscache"
```

```go
rediscahe.SetCacheExpiration(2 * time.Minute)
```