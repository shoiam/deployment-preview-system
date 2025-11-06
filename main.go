package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
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
	// fmt.Printf("Data coming from the request: %v", string(requestData))
	fmt.Printf("JSON form of the request payload: %v\n", payload.Ref)
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
	// http.Serve(newNgClient, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Hello from ngrok-delivered Go app.")
	// }))
	ginClient.RunListener(newNgClient)
}
