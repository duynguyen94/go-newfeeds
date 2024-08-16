package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type respBody struct {
	Message string `json:"message"`
}

func triggerGenNewsfeed(userId int) error {
	newsfeedServiceAddr := os.Getenv("NEWSFEED_SERVICE_ADDR") + "/v1/newsfeeds/" + strconv.Itoa(userId)
	postBody, _ := json.Marshal(map[string]string{})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(newsfeedServiceAddr, "application/json", responseBody)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
