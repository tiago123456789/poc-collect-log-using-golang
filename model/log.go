package model

type Log struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}
