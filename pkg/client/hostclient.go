package client

import (
	tls2 "crypto/tls"
	"errors"
	"fmt"
	"github.com/ip-rw/http2"
	"github.com/ip-rw/http2/fasthttp2"
	tls "github.com/refraction-networking/utls"
	"github.com/rocketlaunchr/go-pool"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

var (
	roller, _ = tls.NewRoller()
	h2client = http2.Client{}
)

type CustomHostClient struct {
	fasthttp.HostClient
	ServerName string
	PoolWrap *pool.ItemWrap
}

func (hc *CustomHostClient) Release() {
	hc.HostClient.CloseIdleConnections()
	rt := 5
	for hc.HostClient.ConnsCount() > 0 && rt > 0 {
		time.Sleep(50*time.Millisecond)
		rt--
	}
	if hc.HostClient.ConnsCount() > 0 {
		hc.PoolWrap.MarkAsInvalid()
	}
	hc.HostClient.Addr = ""
	hc.HostClient.TLSConfig.ServerName = ""
	hc.PoolWrap.Return()
	hc.PoolWrap = nil
}

func GetTransport(hc *CustomHostClient) (fasthttp.TransportFunc, error) {
	if !hc.IsTLS {
		return nil, nil
	}
	uconn, err := roller.Dial("tcp4", hc.Addr, hc.ServerName)
	if err != nil || uconn == nil {
		return nil, err
	}
	alpn := uconn.HandshakeState.ServerHello.AlpnProtocol
	//fmt.Println(alpn)
	switch alpn {
	case "h2":
		c2, _ := http2.NewClient(uconn)
		return fasthttp2.Do(c2), nil
	case "http/1.1", "":
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported ALPN: %v\n", alpn))
	}
}
var cache = tls2.NewLRUClientSessionCache(1000)
var tlsCfg = &tls2.Config{
	InsecureSkipVerify: true,
	ClientSessionCache: cache,
	Renegotiation: tls2.RenegotiateFreelyAsClient,
}

func NewHostClient(addr string, sni string, tls bool) (*CustomHostClient, error) {
	if sni == "" {
		sni, _, _ = net.SplitHostPort(addr)
	}
	h := clientpool.Borrow()
	hc := h.Item.(*CustomHostClient)
	hc.PoolWrap = h
	//hc := &h
	hc.Addr = addr
	hc.ServerName = sni
	hc.IsTLS = tls
	if hc.TLSConfig == nil {
		hc.TLSConfig = tlsCfg.Clone()
	}
	hc.TLSConfig.ServerName = hc.ServerName
	var err error
	hc.Transport, err = GetTransport(hc)
	if err != nil {
		h.MarkAsInvalid()
		hc.Release()
		return nil, err
	}
	return hc, nil
}
