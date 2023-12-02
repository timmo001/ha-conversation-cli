package main

type HomeAssistantPipeline struct {
	ConversationEngine   string      `json:"conversation_engine"`
	ConversationLanguage string      `json:"conversation_language"`
	Language             string      `json:"language"`
	Name                 string      `json:"name"`
	SttEngine            string      `json:"stt_engine"`
	SttLanguage          string      `json:"stt_language"`
	TtsEngine            string      `json:"tts_engine"`
	TtsLanguage          string      `json:"tts_language"`
	TtsVoice             string      `json:"tts_voice"`
	WakeWordEntity       interface{} `json:"wake_word_entity"`
	WakeWordID           interface{} `json:"wake_word_id"`
	ID                   string      `json:"id"`
}

type HomeAssistantPipelinesResult struct {
	Pipelines         []HomeAssistantPipeline `json:"pipelines"`
	PreferredPipeline string                  `json:"preferred_pipeline"`
}

type HomeAssistantPipelineRequestInput struct {
	Text string `json:"text"`
}

type HomeAssistantPipelineRequest struct {
	ID             int                               `json:"id"`
	Type           string                            `json:"type"`
	ConversationID *string                           `json:"conversation_id"`
	StartStage     string                            `json:"start_stage"`
	Input          HomeAssistantPipelineRequestInput `json:"input"`
	EndStage       string                            `json:"end_stage"`
	Pipeline       string                            `json:"pipeline"`
}

type HomeAssistantSpeech struct {
	Speech    string      `json:"speech"`
	ExtraData interface{} `json:"extra_data"`
}

type HomeAssistantSpeechRoot struct {
	Plain HomeAssistantSpeech `json:"plain"`
	SSML  HomeAssistantSpeech `json:"ssml"`
}

type HomeAssistantSuccess struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

type HomeAssistantIntentOutput struct {
	Response struct {
		Speech       HomeAssistantSpeechRoot `json:"speech"`
		Card         map[string]string       `json:"card"`
		Language     string                  `json:"language"`
		ResponseType string                  `json:"response_type"`
		Data         struct {
			Targets []interface{}          `json:"targets"`
			Success []HomeAssistantSuccess `json:"success"`
			Failed  []interface{}          `json:"failed"`
		} `json:"data"`
	} `json:"response"`
	ConversationID interface{} `json:"conversation_id"`
}

type HomeAssistantIntentOutputRoot struct {
	IntentOutput HomeAssistantIntentOutput `json:"intent_output"`
}

type HomeAssistantEvent struct {
	Type      string                        `json:"type"`
	Data      HomeAssistantIntentOutputRoot `json:"data"`
	Timestamp string                        `json:"timestamp"`
}

type HomeAssistantEventRoot struct {
	ID    int                `json:"id"`
	Type  string             `json:"type"`
	Event HomeAssistantEvent `json:"event"`
}
