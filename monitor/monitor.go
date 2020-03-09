package main

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/libp2p/go-libp2p/p2p/protocol/ping"

    "github.com/Multi-Tier-Cloud/common/p2pnode"
)

// collect pings all peers in the monitor node's Peerstore to collect
// performance data. For now it prints out what it finds
func collect(node *p2pnode.Node) {
    // Setup new timer to allow one ping per second
    timer := time.NewTimer(1 * time.Second)
    // Loop infinitely
    for {
        // Get peer in Peerstore
        for _, id := range node.Host.Peerstore().Peers() {
            // Block until timer fires
            <-timer.C
            // Ping and print result
            responseChan := ping.Ping(node.Ctx, node.Host, id)
            result := <-responseChan
            if result.RTT == 0 {
                fmt.Println("ID:", id, "Failed to ping, RTT = 0")
                continue
            }
            fmt.Println("ID:", id, "RTT:", result.RTT)
        }
    }
}

func main() {

    // Setup node
    ctx := context.Background()

    config := p2pnode.NewConfig()
    // Set no rendezvous (anonymous mode)
    config.Rendezvous = []string{}
    node, err := p2pnode.NewNode(ctx, config)
    if err != nil {
        panic(err)
    }

    fmt.Println("Starting data collection")
    collect(&node)
    panic(errors.New("Monitor node exitted monitor loop"))
}
