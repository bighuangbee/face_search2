package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func HttpPost(addr string, data map[string]interface{}) ([]byte, error) {
	jsonData, _ := json.Marshal(&data)
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responeData := []byte{}
	_, err = resp.Body.Read(responeData)
	return responeData, err
}
