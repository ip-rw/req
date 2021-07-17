package client

import (
	tls2 "crypto/tls"
	"errors"
	"fmt"
	"github.com/ip-rw/req/pkg/http2"
	"github.com/ip-rw/req/pkg/http2/fasthttp2"
	tls "github.com/refraction-networking/utls"
	"github.com/rocketlaunchr/go-pool"
	"github.com/valyala/fasthttp"
	"net"
)

var (
	roller, _ = tls.NewRoller()
	h2client  = http2.Client{}
)

type CustomHostClient struct {
	fasthttp.HostClient
	ServerName string
	PoolWrap   *pool.ItemWrap
	Conn       net.Conn
}

func (hc *CustomHostClient) Release() {
	defer func() {
		if err := recover(); err != nil {
			println(err)
		}
	}()
	if hc != nil {
		hc.HostClient.CloseIdleConnections()
		//hc.HostClient.Addr = ""
		//hc.HostClient.TLSConfig.ServerName = ""
		//if hc.Conn != nil {
		//println(hc.Conn.Close())
		//}
		//hc.PoolWrap.Return()
		if hc.Conn != nil {
			hc.Conn.Close()
		}
		hc = nil
		//hc.PoolWrap = nil
	}
}

func GetTransport(hc *CustomHostClient) (fasthttp.TransportFunc, error) {
	if !hc.IsTLS {
		return nil, nil
	}

	uconn, err := roller.Dial("tcp4", hc.Addr, hc.ServerName)
	if err != nil || uconn == nil {
		return nil, err
	}
	hc.Conn = uconn
	alpn := uconn.HandshakeState.ServerHello.AlpnProtocol
	switch alpn {
	case "h2":
		c2, err := http2.NewClient(uconn)
		if err != nil {
			return nil, err
		}
		return fasthttp2.Do(c2), nil
	case "http/1.1", "":
		uconn.Close()
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported ALPN: %v\n", alpn))
	}
}

var cache = tls2.NewLRUClientSessionCache(1000)
var tlsCfg = &tls2.Config{
	InsecureSkipVerify: true,
	ClientSessionCache: cache,
	Renegotiation:      tls2.RenegotiateFreelyAsClient,
}

func NewHostClient(addr string, sni string, tls bool) (*CustomHostClient, error) {
	if sni == "" {
		sni, _, _ = net.SplitHostPort(addr)
	}
	hc := &CustomHostClient{
		HostClient: fasthttp.HostClient{
			Addr:      addr,
			TLSConfig: tlsCfg.Clone(),
		},
	}
	//h := clientpool.Borrow()
	//hc := h.Item.(*CustomHostClient)
	//hc.HostClient.SetMaxConns(1)
	//hc.PoolWrap = h
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
		//h.MarkAsInvalid()
		hc.Release()
		return nil, err
	}
	return hc, nil
}
