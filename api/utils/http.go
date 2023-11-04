package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sankalp-12/clip-url/models"
)

func Put(ctx context.Context, data string) (bool, string) {
	url := "http://localhost:8080/api/v1/shorten"

	body := fmt.Sprintf(`{"url": "%s"}`, string(data))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println("Internal server error: Unable to create shorten request")
		return false, ""
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Internal server error: Unable to send shorten request")
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Internal server error: Shorten request failed")
		return false, ""
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Internal server error: Unable to read shorten response")
		return false, ""
	}

	var jsonResp models.Response
	err = json.Unmarshal(response, &jsonResp)
	if err != nil {
		log.Println("Internal server error: Unable to parse shorten response")
		return false, ""
	}

	return true, jsonResp.NewURL
}
