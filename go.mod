module github.com/hkmadao/rtmp2flv

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	github.com/deepch/vdk v0.0.0-20210523103705-5b25bda1a000
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
)

replace github.com/deepch/vdk => github.com/hkmadao/vdk v0.0.0-20241120073805-439b6309323c
