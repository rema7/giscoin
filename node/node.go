package node

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	log "github.com/inconshreveable/log15"
	"os"
)

const messageId = 0

type Message string

var (
	proto = p2p.Protocol{
		Name:    "ping",
		Version: 1,
		Length:  1,
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
			message := "ping"

			// Sending the message to connected peer
			err := p2p.Send(rw, 0, message)
			if err != nil {
				return fmt.Errorf("Send message fail: %v", err)
			}
			fmt.Println("sending message", message)

			// Receiving the message from connected peer
			received, err := rw.ReadMsg()
			if err != nil {
				return fmt.Errorf("Receive message fail: %v", err)
			}

			var myMessage string
			err = received.Decode(&myMessage)

			fmt.Println("received message", string(myMessage))

			return nil
		},
	}
)

//func MyProtocol() p2p.Protocol {
//	return p2p.Protocol{
//		Name:    "MyProtocol",
//		Version: 1,
//		Length:  1,
//		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
//			message := "ping"
//
//			// Sending the message to connected peer
//			err := p2p.Send(rw, 0, message)
//			if err != nil {
//				return fmt.Errorf("Send message fail: %v", err)
//			}
//			fmt.Println("sending message", message)
//
//			// Receiving the message from connected peer
//			received, err := rw.ReadMsg()
//			if err != nil {
//				return fmt.Errorf("Receive message fail: %v", err)
//			}
//
//			var myMessage string
//			err = received.Decode(&myMessage)
//
//			fmt.Println("received message", string(myMessage))
//
//			return nil
//		},
//	}
//}

func Start(port int) *p2p.Server {
	privateKey, _ := crypto.GenerateKey()

	config := p2p.Config{
		MaxPeers:        10,
		PrivateKey:      privateKey,
		Name:            "my node name",
		ListenAddr:      fmt.Sprintf(":%d", port),
		Protocols:       []p2p.Protocol{proto},
		EnableMsgEvents: true,
	}
	srv := &p2p.Server{
		Config: config,
	}

	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nodeInfo := srv.NodeInfo()
	log.Info("server started", "enode", nodeInfo.Enode, "name", nodeInfo.Name, "ID", nodeInfo.ID, "IP", nodeInfo.IP)

	return srv
}

func ConnectToPeer(srv *p2p.Server, enode string) error {
	node, err := discover.ParseNode(enode)
	if err != nil {
		log.Error("Bootstrap URL invalid", "enode", enode, "err", err)
		return err
	}
	srv.AddPeer(node)

	log.Info("after add", "node one peers", srv.Peers())
	return nil
}

func SubscribeToEvents(srv *p2p.Server, communicated chan<- bool) {
	// Subscribing to the peer events
	peerEvent := make(chan *p2p.PeerEvent)
	eventSub := srv.SubscribeEvents(peerEvent)

	for {
		select {
		case event := <-peerEvent:
			if event.Type == p2p.PeerEventTypeMsgRecv {
				log.Info("Received message received notification")
				communicated <- true
			}
		case <-eventSub.Err():
			log.Info("subscription closed")

			// Closing the channel so that server gets stopped since
			// there won't be any more events coming in
			close(communicated)
		}
	}
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
