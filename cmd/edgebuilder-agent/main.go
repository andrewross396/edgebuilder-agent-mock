package main

import (
	"edgebuilder-agent-mock/agent"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var payload agent.AgentRequest
	json.Unmarshal(msg.Payload(), &payload)
	time.Sleep(10 * time.Second)
	fmt.Printf("\nReceived operation: %s from topic: %s\n", payload.Ops[0].Name, msg.Topic())

	ops := agent.AgentResponseOp{
		ID:            payload.Ops[0].ID,
		Success:       true,
		Result:        nil,
		Errors:        nil,
		StartTime:     12345,
		ExecutionTime: 12345,
	}

	var response agent.AgentResponse
	response.Version = payload.Version
	response.Signature = "test"
	response.ResponderID = strings.Split(msg.Topic(), "/")[3]
	response.RequestID = payload.RequestID
	response.Success = true
	response.Ops = []agent.AgentResponseOp{ops}

	responsePayloadJson, _ := json.Marshal(&response)
	responsePayload := string(responsePayloadJson)

	topic := fmt.Sprintf("/eb/command/%s", payload.RequestID)
	token := client.Publish(topic, 0, false, responsePayload)
	token.Wait()
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	keepAlive := make(chan os.Signal)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	var broker = "127.0.0.1"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("agent")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)

	<-keepAlive
}

func sub(client mqtt.Client) {
	topic := "/eb/command/123e4567-e89b-12d3-a456-426614174001"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)

	topic2 := "/eb/command/6a55a0a2-68f7-4964-8e06-607e7fadcf89"
	token2 := client.Subscribe(topic2, 1, nil)
	token2.Wait()
	fmt.Printf("Subscribed to topic: %s", topic2)

}
