package client

import (
	"github.com/rocketlaunchr/go-pool"
	"github.com/valyala/fasthttp"
)

var clientpool = NewHostClientPool(10000)
func NewHostClientPool(max int) pool.Pool {
	pool := pool.New(max) // maximum of 5 items in pool
	pool.SetFactory(func() interface{} {
		return &CustomHostClient{
			HostClient: fasthttp.HostClient{},
		}
	})
	return pool
}

