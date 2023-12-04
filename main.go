package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

const (
	AppVersion = "0.1.0"
	configFile = "config.json"
)

var (
	config                  Settings
	dialer                  = websocket.DefaultDialer
	homeAssistantConnection *websocket.Conn
	homeAssistantPipelines  HomeAssistantResponse[HomeAssistantPipelinesResult]
)

// Load config file
func loadConfig() Settings {
	// 1. Check if config file exists
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		// If not, create it
		log.Print("Config file not found, creating it...")
		file, err := os.Create("config.json")
		if err != nil {
			log.Fatal("Error creating config file:", err)
		}
		defer file.Close()

		// Write default config to file
		defaultConfig := Settings{
			HomeAssistantHost:  "homeassistant.local",
			HomeAssistantPort:  8123,
			HomeAssistantToken: "",
		}
		encoder := json.NewEncoder(file)
		err = encoder.Encode(defaultConfig)
		if err != nil {
			log.Fatal("Error writing default config to file:", err)
		}
	}

	// 2. Open config file
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer file.Close()

	// 3. Read config file to settings struct
	decoder := json.NewDecoder(file)
	newConfig := Settings{}
	err = decoder.Decode(&newConfig)
	if err != nil {
		log.Fatal("Error reading config file. Error was: ", err)
	}

	// 4. Check if config is empty
	if newConfig.HomeAssistantHost == "" || newConfig.HomeAssistantPort == 0 || newConfig.HomeAssistantToken == "" {
		log.Fatal("Config file is empty, please fill it with your Home Assistant details")
	}

	// 5. Return config
	return newConfig
}

// Setup connection to Home Assistant via WebSocket
func setupHomeAssistantConnection() *websocket.Conn {
	// 1. Setup WebSocket URL
	wsUrl := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", config.HomeAssistantHost, config.HomeAssistantPort), Path: "/api/websocket"}
	log.Print("Connecting to Home Assistant at: ", wsUrl.String())

	// 2. Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}

	// 3. Read welcome message
	var welcomeResponse HomeAssistantResponse[any]
	err = conn.ReadJSON(&welcomeResponse)
	if err != nil {
		log.Fatal("Error reading welcome message:", err)
	}
	log.Print("First response: ", welcomeResponse)

	// 4. Create auth message
	authMessage := HomeAssistantAuthRequest{
		Type:        "auth",
		AccessToken: config.HomeAssistantToken,
	}

	// 5. Send auth message
	err = conn.WriteJSON(authMessage)
	if err != nil {
		log.Fatal("Error sending auth message:", err)
	}

	// 6. Read auth response
	var authResponse HomeAssistantResponse[any]
	err = conn.ReadJSON(&authResponse)
	if err != nil {
		log.Fatal("Error reading auth response:", err)
	}
	log.Print("Auth response: ", authResponse)

	return conn
}

func homeAssistantGetPipelines() HomeAssistantResponse[HomeAssistantPipelinesResult] {
	// Create message
	message := HomeAssistantRequest{
		ID:   1,
		Type: "assist_pipeline/pipeline/list",
	}

	// Send message
	err := homeAssistantConnection.WriteJSON(message)
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	// Read response
	var pipelinesResponse HomeAssistantResponse[HomeAssistantPipelinesResult]
	err = homeAssistantConnection.ReadJSON(&pipelinesResponse)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}
	log.Print("List pipelines response: ", pipelinesResponse)

	return pipelinesResponse
}

func main() {
	// Create file logger
	logFile, err := os.OpenFile("ha-conversation-cli.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Print("HA Conversation CLI: ", AppVersion)

	// Load config
	config = loadConfig()

	// Setup connection to Home Assistant
	homeAssistantConnection = setupHomeAssistantConnection()

	// Get pipelines
	homeAssistantPipelines = homeAssistantGetPipelines()

	requestInputText := flag.String("text", "", "Text to send to Home Assistant")
	flag.Parse()

	if *requestInputText == "" {
		log.Fatal("No text provided")
	}

	// Create message
	message := HomeAssistantPipelineRequest{
		ID:             2,
		Type:           "assist_pipeline/run",
		ConversationID: nil,
		StartStage:     "intent",
		Input: HomeAssistantPipelineRequestInput{
			Text: *requestInputText,
		},
		EndStage: "intent",
		Pipeline: homeAssistantPipelines.Result.PreferredPipeline,
	}
	// Log message as JSON
	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Error marshalling message:", err)
	}
	log.Print("Sending message: ", string(messageJson))

	// Send message
	err = homeAssistantConnection.WriteJSON(message)
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	// Read result
	var result HomeAssistantResponse[any]
	err = homeAssistantConnection.ReadJSON(&result)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}
	log.Print("Response: ", result)

	var event HomeAssistantEventRoot

	// Read next message until we get type event and data.type is intent-end
	for event.Type != "event" || event.Event.Type != "intent-end" {
		err = homeAssistantConnection.ReadJSON(&event)
		if err != nil {
			log.Fatal("Error reading response:", err)
		}
		log.Print("Response: ", event)
	}

	log.Print("Response: ", event.Event.Data.IntentOutput.Response.Speech.Plain.Speech)

	// Close connection
	homeAssistantConnection.Close()

	log.Print("Done")

	fmt.Println(event.Event.Data.IntentOutput.Response.Speech.Plain.Speech)
}
