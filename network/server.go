package network

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	log "github.com/inconshreveable/log15"
	"os"
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
