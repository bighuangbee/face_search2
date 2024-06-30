package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {

	url := "http://localhost:6002/face/unregiste/all"
	data := `{}`
	payload := bytes.NewBufferString(data)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusOK {

	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(string(respBody))
	}

}
