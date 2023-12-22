package main

import (
"time"
"go.mongodb.org/mongo-driver/bson"
"log"
"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	
)
//mongodb+srv://admin1:admin123@cluster0.4df1svs.mongodb.net/?retryWrites=true&w=majority
const help = `...`
var Hoster host.Host
var mongoDB *MongoDB


func CheckNodesPeriodically() {
	// Run the check periodically every 5 minutes
	ticker := time.NewTicker(60* time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking nodes...")

			// Initialize MongoDB
			
			

			// Read all documents from the peermeta collection
			var nodes []NodeInfo
			err := mongoDB.GetDocuments(PeerMeta, bson.D{}, &nodes)
			if err != nil {
				log.Fatal(err)
			}

			// Ping each node and remove it if it's down
			for _, node := range nodes {
				
				if err := pingNode(node.ID); err != nil {
					fmt.Printf("Node %s is down. Removing from the database.\n", node.ID)
					err := mongoDB.DeleteDocument(PeerMeta, node.ID)
					if err != nil {
						log.Printf("Error removing node %s from the database: %v\n", node.URL, err)
					}
				}
			}
		}
	}
}
func pingNode(peerID string) error {
	// Create a new ping service
	
	// Ping the specified peer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println("the peer id is ", peerID)
	defer cancel()
		peerid,err := peer.Decode(peerID)
		if err != nil{
			return err
		}
		pingResults := ping.Ping(ctx, Hoster, peerid)
	
		// Wait for the result or context cancellation
		select {
		case result := <-pingResults:
			// Check if the ping was successful (got an RTT)
			fmt.Println(result)
			if result.Error == nil {
				RequestMeta(peerid)
				return nil
			}
			// Ping failed
			return result.Error
	
		case <-ctx.Done():
			// Context canceled
			return ctx.Err()
		}
	
}
func main() {
	flag.Usage = func() {
		fmt.Println(help)
		flag.PrintDefaults()
	}
	destPeer := flag.String("d", "", "destination peer address")
	port := flag.Int("p", 9900, "proxy port")
	p2pport := flag.Int("l", 12000, "libp2p listen port")
	
	mongourl := flag.String("m", "", "mongodb url address")
	flag.Parse()
	var err error
	fmt.Println(*mongourl)
	
	mongoDB, err = NewMongoDB(*mongourl, "node", NodeMeta, PeerMeta)
if err != nil {
	fmt.Println("the error is in connection")
    fmt.Println(err)
}
err = mongoDB.DropAllCollections()
if err != nil {
	fmt.Println("the error in dropping connection")
    fmt.Println(err)
}
fmt.Println(mongoDB)
// Add a new collection dynamically

/**node := NodeInfo{
	ID:  "some_id",
	URL: "https://example.com",
	// Add values for additional fields if needed
}

// Insert the document into the collection
err = mongoDB.InsertDocument("Trial", node)**/
if err != nil {
	fmt.Println(err)
}
	

	

	if *destPeer != "" {
		Hoster = makeRandomHost(*p2pport)
		Hoster.SetStreamHandler(Protocol, StreamHandler)
		go CheckNodesPeriodically()
		Serve(*port)
	} else {
		Hoster = makeRandomHost(*p2pport)
		Hoster.SetStreamHandler(Protocol, StreamHandler)
		for _, a := range Hoster.Addrs() {
			fmt.Printf("%s/ipfs/%s\n", a, Hoster.ID())
		}
		go CheckNodesPeriodically()
		Serve(*port)
	}
}
