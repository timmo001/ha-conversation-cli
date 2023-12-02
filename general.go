package main

type Settings struct {
	HomeAssistantHost  string
	HomeAssistantPort  int
	HomeAssistantToken string
}

type HomeAssistantRequest struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type HomeAssistantAuthRequest struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type HomeAssistantResponse[T any] struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result    T      `json:"result"`
	HaVersion string `json:"ha_version"`
}
