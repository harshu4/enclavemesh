package main

import (
	
	"encoding/json"
	"fmt"
	
	"log"
	"github.com/libp2p/go-libp2p/core/host"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"math/rand"
	"strconv"
	ma "github.com/multiformats/go-multiaddr"
)

type Node struct {
	ID     string `json:"id"`
	PeerID peer.ID `json:"peerid"`
}

type Url struct {
	Description string `json:"description"`
	Url string `json:"url"`
	Seconds string `json:"interval"`
}

type Removal struct {
	Id int `json:"id"`
	
}

type AddCollection struct {
	Id int `json:"id"`
	Peerid peer.ID `json:"peerid"`
}

var nodes []Node

func AddAddrToPeerstore(h host.Host, addr string) peer.ID {
	ipfsaddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		log.Fatalln(err)
	}
	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peerid))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	h.Peerstore().AddAddr(peerid, targetAddr, peerstore.PermanentAddrTTL)
	return peerid
}


func AddNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	var newNode Node
	err := json.NewDecoder(r.Body).Decode(&newNode)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	destPeerID := AddAddrToPeerstore(Hoster, newNode.ID)
	
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newNode.PeerID = destPeerID
	fmt.Println("hello how ")
	fmt.Println(destPeerID)
	
	nodes = append(nodes, newNode)
	node := NodeInfo{
		ID:  destPeerID.String(),
		URL: newNode.ID,
		// Add values for additional fields if needed
	}
	
	// Insert the document into the collection
	fmt.Println("where is the problem")
	err = mongoDB.InsertDocument(PeerMeta, node)
	if err != nil {
		fmt.Println("Failed to add node to db")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(newNode)

}

func GetNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func AddData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "POST"{
	var tee Url;
	random := rand.Intn(9999999)
	err := json.NewDecoder(r.Body).Decode(&tee)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.Body)

	fmt.Println(tee.Url)
	urldata := MetaInfo{

		DataID: random,
		Title: tee.Url,
		Description: tee.Description,
		Seconds: tee.Seconds,
		Working: true,
	}
	json.NewEncoder(w).Encode(urldata)
	mongoDB.InsertDocument(NodeMeta,urldata)
	secondint,err := strconv.Atoi(tee.Seconds)
	if err != nil {
		fmt.Println(err)
	}

	go StartWork(random,tee.Url, secondint)
}

}
func convertMtoD(m bson.M) bson.D {
    var d bson.D
    for key, value := range m {
        d = append(d, bson.E{Key: key, Value: value})
    }
    return d
}

func RemoveData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "POST"{
	var stopid Removal;
	err := json.NewDecoder(r.Body).Decode(&stopid)
	if err != nil {
		fmt.Println(err)
	}
	
	filter := bson.M{"data_id": stopid.Id}
	
	
	mongoDB.GetDocumentAndEdit(NodeMeta,filter)
	json.NewEncoder(w).Encode(stopid)
	StopWork(stopid.Id)

}
}

func ReqCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "POST"{
	
	var collectparam AddCollection;
	err := json.NewDecoder(r.Body).Decode(&collectparam)
	if err != nil {
		fmt.Println(err)
	}
	RequestCollections(collectparam.Peerid,collectparam.Id,)
	json.NewEncoder(w).Encode(collectparam.Id)
}
}

func ReqData(w http.ResponseWriter, r *http.Request) {
	var dat []MetaInfo
	filter := bson.D{{}}
	mongoDB.GetDocuments(NodeMeta,filter,&dat)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dat)
}


func ReqPeerData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var newNode Node
	json.NewDecoder(r.Body).Decode(&newNode)
	var dat []MetaInfo
	fmt.Println(dat)
	filter := bson.D{{}}
	fmt.Println(PeerMeta+newNode.ID)
	mongoDB.GetDocuments(PeerMeta+newNode.ID,filter,&dat)
	fmt.Println(dat)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dat)
}

