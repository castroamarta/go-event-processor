package types

type Event struct {
	Timestamp      string `json:"timestamp"`
	TranslationID  string `json:"translation_id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	ClientName     string `json:"client_name"`
	EventName      string `json:"event_name"`
	NrWords        int    `json:"nr_words"`
	Duration       int    `json:"duration"`
}

type DeliveryTimes struct {
	DurationSum int
	Count       int
}

type AverageDeliveryTimes struct {
	Timestamp string
	Average   float64
}
