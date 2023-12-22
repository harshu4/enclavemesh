package main

import (
	
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"context"
	"go.mongodb.org/mongo-driver/bson"
"strconv"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"reflect"
	
	
)

const Protocol = "/proxy-example/0.0.1"

type MessageType string

const (
	PingPongMessageType MessageType = "pingpong"
	RequestMetaType MessageType = "requestmeta"
	RespondMeta MessageType = "respondmeta"
	RequestCollection MessageType = "requestcollection"
	RespondCollection MessageType = "respondcollection"

)

type Request struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type Response struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data,omitempty"`
}


func convertToInterfaceSlice(data interface{}) ([]interface{}, error) {
	// Check if data is a slice
	sliceValue := reflect.ValueOf(data)
	if sliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice")
	}

	// Create a new slice of interfaces
	result := make([]interface{}, sliceValue.Len())

	// Copy elements from the original slice to the new slice
	for i := 0; i < sliceValue.Len(); i++ {
		result[i] = sliceValue.Index(i).Interface()
	}

	return result, nil
}


func RequestMeta(destPeerID peer.ID) error {
	stream, err := Hoster.NewStream(context.Background(), destPeerID, Protocol)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer stream.Close()

	response := Response{Type: RequestMetaType, Data: "pong"}
	err = json.NewEncoder(stream).Encode(response)
	if err != nil{
		return err
	}
	StreamHandler(stream)
	return nil
}


func RequestCollections(destPeerID peer.ID,collectionid int) error {
	stream, err := Hoster.NewStream(context.Background(), destPeerID, Protocol)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer stream.Close()

	response := Response{Type: RequestCollection, Data: collectionid}
	err = json.NewEncoder(stream).Encode(response)
	if err != nil{
		return err
	}
	StreamHandler(stream)
	return nil
}

func StreamHandler(s network.Stream) {
	defer s.Close()
	fmt.Printf("Received a stream from %s\n", s.Conn().RemotePeer())

	var request Request
	err := json.NewDecoder(s).Decode(&request)
	if err != nil {
		log.Printf("Error decoding request: %v", err)
		return
	}

	if request.Type == PingPongMessageType {
		fmt.Println(request.Type)
		response := Response{Type: PingPongMessageType, Data: "pong"}

		err := json.NewEncoder(s).Encode(response)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			return
		}
	}

	if request.Type == RequestMetaType {
		data := []MetaInfo{}
		filter := bson.D{{}}
		err := mongoDB.GetDocuments(NodeMeta,filter,&data)
		if err != nil{
			log.Println(err)
			log.Printf("error in getting documents")
		}
		response := Response{Type: RespondMeta,Data:data}
		err = json.NewEncoder(s).Encode(response)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			return
		}
	}
	if request.Type == RequestCollection {
		data := []Data{}
		filter := bson.D{{}}
		collname := CollectionPrefix+ strconv.FormatFloat(request.Data.(float64), 'f', -1, 64)
		err := mongoDB.GetDocuments(collname,filter,&data)
		if err != nil{
			log.Println(err)
			log.Printf("error in getting documents")
		}
		response := Response{Type: RespondCollection,Data:data}
		err = json.NewEncoder(s).Encode(response)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			return
		}
	}
	if request.Type == RespondMeta {

	docs,err :=  convertToInterfaceSlice(request.Data)
	if err != nil {
		log.Println(err)
	}
	
	
		collname := PeerMeta+ s.Conn().RemotePeer().String()
		err = mongoDB.DropCollection(collname)
		err = mongoDB.AddCollection(collname)
		if err != nil {
			log.Println(err)
		}
		err = mongoDB.InsertManyDocuments(collname,docs)
		if err != nil {
			log.Println(err)
		}
	
}else {
	fmt.Println("docs is empty")
}
	if request.Type == RespondCollection {

		docs,err :=  convertToInterfaceSlice(request.Data)
		fmt.Println("why would this work ok ok ")
		if err != nil {
			log.Println(err)
		}
		var dataid interface{};
	dataSlice, _ := request.Data.([]interface{});
	fmt.Println(reflect.TypeOf(request.Data))
	if len(dataSlice) > 0 {
		// Check if the underlying type is Document
		fmt.Println(reflect.TypeOf(dataSlice[0]))
		
		if doc, ok := dataSlice[0].(map[string]interface{}); ok {
			// Access the DataID field
			
			 dataid = doc["DataID"]
		} else {
			fmt.Println("docs[0] is not of type Document")
		}
			collname := PeerCollection+strconv.FormatFloat(dataid.(float64), 'f', -1, 64)
			err = mongoDB.DropCollection(collname)
			err = mongoDB.AddCollection(collname)
			if err != nil {
				log.Println(err)
			}
			err = mongoDB.InsertManyDocuments(collname,docs)
			if err != nil {
				log.Println(err)
			}
		}
}}

func makeRandomHost(port int) host.Host {
	host, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)))
	if err != nil {
		log.Fatalln(err)
	}
	return host
}

func Serve(port int) {
	fmt.Println("is this updating")
	fmt.Printf("hello from file")
	http.HandleFunc("/nodes/data", AddData)
	http.HandleFunc("/nodes/rdata", RemoveData)
	http.HandleFunc("/nodes/add", AddNode)
	http.HandleFunc("/nodes", GetNode)
	http.HandleFunc("/collection/req", ReqCollection)
	http.HandleFunc("/nodes/getdata", ReqData)
	http.HandleFunc("/nodes/peerdata", ReqPeerData)
	
	
	fmt.Printf("Server running on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
