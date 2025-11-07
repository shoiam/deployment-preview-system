package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shoiam/deployment-preview-system.git/dockerClient"
	"golang.ngrok.com/ngrok/v2"
)

type WebhookPayload struct {
	Ref string `json:"ref"`
}

func handleWebhook(c *gin.Context) {
	fmt.Println("I am in the handleWebhook")
	data := c.Request.Body
	requestData, err := io.ReadAll(data)
	var payload WebhookPayload
	json.Unmarshal(requestData, &payload)

	if err != nil {
		log.Fatal("Unable to read request data.")
	}
	fmt.Printf("JSON form of the request payload: %v\n", payload.Ref)
	dockerClient.ClientElement()
}

func main() {
	ginClient := gin.Default()
	ctx := context.Background()
	newNgClient, err := ngrok.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App public URL is: %v\n", newNgClient.URL())

	ginClient.POST("/webhook", handleWebhook)
	ginClient.RunListener(newNgClient)
}
