package main

//forwards http/https requests in a fire and forget way
//uses filesystem for persistence

import (
	"fmt"
	"net/http"
	"strconv"
)

const NWORKERS = 1

func forwardRequest(client *http.Client, reqChan chan *http.Request) {
	for r := range reqChan {
		if len(r.Header["Forwardhost"]) >= 1 {
			host := r.Header["Forwardhost"][0]
			myReq, _ := http.NewRequest(r.Method, "http://"+host+r.URL.Path, r.Body)
			myReq.Header = r.Header
			n, err := strconv.Atoi(r.URL.Path[1:])
			fmt.Println(n)
			if err!=nil {
				fmt.Println(err)
			}
			myReq.URL.Path =  "/"+strconv.Itoa(n + 1)
			//delete(myReq.Header, "Forwardhost")
			data, err := client.Do(myReq)
			data.Body.Close()
			if err!=nil {
				fmt.Println(err)
			}
			fmt.Println(data)
		} else {
			fmt.Println("No Header found. Nothing to do")
		}
	}

}

func handleRequest(reqChan chan *http.Request) func(w http.ResponseWriter, r *http.Request) {
	//create a function for HandleRequest with a channel in scope
	return func(w http.ResponseWriter, r *http.Request) {
		//returns a 202 and pushes request toa channel
		select {
		case reqChan <- r:
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Working on it!::"+r.URL.Path))
			fmt.Println("request handled")
		default:
			w.WriteHeader(http.StatusGatewayTimeout)
			fmt.Println("we are full!")

		}
	}
}

func main() {
	reqChan := make(chan *http.Request, NWORKERS)
	//create an http server
	fmt.Println("Entering...")
	//TODO channel for request + workers
    
	client := &http.Client{}
	for i := 0; i < NWORKERS; i++ {
		go forwardRequest(client, reqChan)
	}

	http.HandleFunc("/", handleRequest(reqChan))
	fmt.Println(http.ListenAndServe(":8080", nil))
}
