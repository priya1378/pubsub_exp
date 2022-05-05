package main

import (
	"fmt"
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "pubsub-chat-example"

func main() {
	ctx := context.Background()
	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	// setup local mDNS discovery
	if err := setupDiscovery(h); err != nil {
		panic(err)
	}

	//ps.Join("hey")
	
	//topic, err := ps.Join("hey")
	PubSub(ctx, ps, h.ID())
	if err != nil {
		panic(err)
	}
	//sub, err := topic.Subscribe()*/
	
}

// discoveryNotifee g	ets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

func PubSub(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID){
	topic, err := ps.Join("hey")
	if err != nil {
		panic(err)
	}

	sub, err := topic.Subscribe()

	if err!= nil {
		panic(err)
	}
	msg, err := sub.Next(ctx)
	if err!=nil{
		panic(err)
	}
	fmt.Println(msg)
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

