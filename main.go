package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const baseAddr = "http://localhost:8080"

func main() {
	httpTest(
		"look for 404",
		"GET",
		fmt.Sprintf("%s/", baseAddr),
		nil,
		nil,
	)
}

func httpTest(
	description string,
	method,
	url string,
	reqStruct interface{},
	respStruct interface{},
) {
	fmt.Printf("======== %v ========\n", description)
	defer fmt.Printf("======== END ========\n")

	var bodyReader io.Reader
	if reqStruct != nil {
		dat, err := json.Marshal(reqStruct)
		if err != nil {
			log.Printf("json.Marshal body: %v\n", err)
			return
		}
		bodyReader = bytes.NewReader(dat)
	}

	fmt.Printf("Sending %s request to %s\n", method, url)
	cacheBuster := rand.Int()
	req, err := http.NewRequest(method, fmt.Sprintf("%s?v=%v", url, cacheBuster), bodyReader)
	if err != nil {
		fmt.Printf("http.NewRequest: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("http.Do: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response received!")
	fmt.Println("Status code:", resp.StatusCode)
	if resp.StatusCode > 299 {
		return
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("io.ReadAll: %v\n", err)
		return
	}

	if respStruct == nil {
		fmt.Printf("Response body:\n%s\n", string(dat))
		return
	}

	err = json.Unmarshal(dat, respStruct)
	if err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		return
	}
	parsedDat, err := json.MarshalIndent(respStruct, "", "  ")
	if err != nil {
		log.Printf("json.MarshalIndent: %v\n", err)
		return
	}
	fmt.Printf("Parsed resp body: %s\n", string(parsedDat))
}
