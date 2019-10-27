## Migrate from glide to go modules 

[ref](https://blog.liquidbytes.net/2018/09/quick-and-easy-guide-for-migrating-to-go-1-11-modules/)

```bash
go build ./...
env GO111MODULE=on go mod init
env GO111MODULE=on go mod vendor
```

## Ignore permission issues

go get https://github.com/martin-helmich/cloudnativego-backend

glide get github.com/mitchellh/mapstructure
glide get github.com/martin-helmich/cloudnativego-backend
glide get github.com/martin-helmich/cloudnativego-backend/src/lib/helper/kafka
glide get github.com/martin-helmich/cloudnativego-backend/src/lib/msgqueue
