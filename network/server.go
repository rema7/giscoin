package network

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	log "github.com/inconshreveable/log15"
	"os"
	"time"
)

const messageId = 0

type Message string

func MyProtocol() p2p.Protocol {
	return p2p.Protocol{
		Name:    "MyProtocol",
		Version: 1,
		Length:  1,
		Run:     msgHandler,
	}
}

func Start() {
	nodekey, _ := crypto.GenerateKey()
	config := p2p.Config{
		MaxPeers:   10,
		PrivateKey: nodekey,
		Name:       "my node name",
		ListenAddr: ":30300",
		Protocols:  []p2p.Protocol{MyProtocol()},
	}
	srv := p2p.Server{
		Config: config,
	}

	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nodeInfo := srv.NodeInfo()
	log.Info("server started", "enode", nodeInfo.Enode, "name", nodeInfo.Name, "ID", nodeInfo.ID, "IP", nodeInfo.IP)
	url := "enode://528ef3a7de7746ccfd5a02f98e5485961974ffa9e23cb9afab155d54a489becf48d97eadb49e75d9c1c3ac27fd6141ccc7e2d87c42b05c09a77c26a32c7c94ce@127.0.0.1:30303"
	node, err := discover.ParseNode(url)
	if err != nil {
		log.Error("Bootstrap URL invalid", "enode", url, "err", err)
	}
	fmt.Printf("%s\n", srv.Self().String())
	fmt.Printf("%s\n", node.String())

	srv.AddPeer(node)
	time.Sleep(time.Millisecond * 100)

	log.Info("after add", "node one peers", srv.Peers())

	select {}
}

func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
	for {
		msg, err := ws.ReadMsg()
		if err != nil {
			return err
		}

		var myMessage Message
		err = msg.Decode(&myMessage)
		if err != nil {
			// handle decode error
			continue
		}

		switch myMessage {
		case "foo":
			err := p2p.SendItems(ws, messageId, "bar")
			if err != nil {
				return err
			}
		default:
			fmt.Println("recv:", myMessage)
		}
	}

	return nil
}
