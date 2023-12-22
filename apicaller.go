package main

import (
	"fmt"
	"sync"
	"time"
	
	"io/ioutil"
	"net/http"
	"strconv"
)

type APICaller struct {
	ID       int
	APIURL   string
	Interval time.Duration
	Stop     chan struct{}
}
var ApiCallers = make(map[int]*APICaller)
func makeGETRequest(url string) (string, error) {
	// Make GET request
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func NewAPICaller(id int, apiURL string, interval time.Duration) *APICaller {
	return &APICaller{
		ID:       id,
		APIURL:   apiURL,
		Interval: interval,
		Stop:     make(chan struct{}),
	}
}

func (ac *APICaller) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(ac.Interval)
	ticker := time.NewTicker(ac.Interval)

	for {
		select {
		case <-ticker.C:
			// Call the API here
			result, err := makeGETRequest(ac.APIURL)
			if err != nil {
				fmt.Println("error fetching url")
			}
			data := Data{
				DataID:ac.ID,
				JSONres: result,
				Signature: "sdd",
				Timestamp:1,
			}
			mongoDB.AddCollection(CollectionPrefix+strconv.Itoa(ac.ID))
			mongoDB.InsertDocument(CollectionPrefix+strconv.Itoa(ac.ID),data)
			fmt.Printf("Calling API %d: %s\n", ac.ID, ac.APIURL)
		case <-ac.Stop:
			fmt.Printf("Stopping API caller %d for: %s\n", ac.ID, ac.APIURL)
			return
		}
	}
}

func (ac *APICaller) StopCaller() {
	close(ac.Stop)
}


func StartWork(random int ,url string,seconds int){
	api1Caller := NewAPICaller(random, url, time.Duration(seconds)*time.Second);
	ApiCallers[api1Caller.ID] = api1Caller
	var wg sync.WaitGroup
	wg.Add(1)
	go api1Caller.Start(&wg)
	wg.Wait()
}




func StopWork(id int){
	ApiCallers[id].StopCaller()
}
