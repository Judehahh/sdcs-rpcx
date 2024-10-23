package rpc

import (
	"fmt"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

func NewXClient(port int) client.XClient {
	addr := fmt.Sprintf("localhost:%d", port)
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	opt := client.DefaultOption
	opt.SerializeType = protocol.JSON

	return client.NewXClient("Cache", client.Failtry, client.RandomSelect, d, opt)
}
