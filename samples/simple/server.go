package main

import (
	"encoding/json"
	"time"

	tp "github.com/henrylee2cn/teleport"
)

func main() {
	go tp.GraceSignal()
	// tp.SetReadLimit(10)
	tp.SetShutdown(time.Second*20, nil, nil)
	var cfg = &tp.PeerConfig{
		DefaultReadTimeout:   time.Minute * 5,
		DefaultWriteTimeout:  time.Millisecond * 500,
		TlsCertFile:          "",
		TlsKeyFile:           "",
		SlowCometDuration:    time.Millisecond * 500,
		DefaultHeaderCodec:   "protobuf",
		DefaultBodyCodec:     "json",
		DefaultBodyGzipLevel: 5,
		PrintBody:            true,
		ListenAddrs: []string{
			"0.0.0.0:9090",
			"0.0.0.0:9091",
		},
	}
	var peer = tp.NewPeer(cfg)
	{
		group := peer.PullRouter.Group("group")
		group.Reg(new(Home))
	}
	peer.PullRouter.SetUnknown(UnknownPullHandle)
	peer.Listen()
}

// Home controller
type Home struct {
	tp.PullCtx
}

// Test handler
func (h *Home) Test(args *map[string]interface{}) (map[string]interface{}, tp.Xerror) {
	h.Session().Push("/push/test?tag=from home-test", map[string]interface{}{
		"your_id": h.Query().Get("peer_id"),
		"a":       1,
	})
	return map[string]interface{}{
		"your_args":   *args,
		"server_time": time.Now(),
	}, nil
}

func UnknownPullHandle(ctx tp.UnknownPullCtx) (interface{}, tp.Xerror) {
	time.Sleep(1)
	var v = struct {
		ConnPort int
		json.RawMessage
		Bytes []byte
	}{}
	codecName, err := ctx.Bind(&v)
	if err != nil {
		return nil, tp.NewXerror(1, err.Error())
	}
	tp.Debugf("UnknownPullHandle: codec: %s, conn_port: %d, RawMessage: %s, bytes: %s",
		codecName, v.ConnPort, v.RawMessage, v.Bytes,
	)
	return []string{"a", "aa", "aaa"}, nil
}
