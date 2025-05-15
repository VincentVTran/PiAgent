package config

type PiOperation struct {
	ClientID   string                 `json:"client_id"`
	Operation  string                 `json:"operation"`
	Parameters map[string]interface{} `json:"parameters"`
}
