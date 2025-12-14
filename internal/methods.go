package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/GrandTheBest/heartbeat/types"
)

func CallG(url string) (*types.StatusResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bodyData types.StatusResponse
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		return nil, err
	}

	return &bodyData, nil
}

func CallP(url string, payload any) (*types.ResponsePacket, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bodyData types.ResponsePacket
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		return nil, err
	}

	return &bodyData, nil
}
