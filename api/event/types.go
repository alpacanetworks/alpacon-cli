package event

import "time"

type EventAttributes struct {
	Server      string `json:"server"`
	Shell       string `json:"shell"`
	Command     string `json:"command"`
	Result      string `json:"result"`
	Status      bool   `json:"status"`
	Operator    string `json:"operator"`
	RequestedAt string `json:"requested_at"`
}

type EventListResponse struct {
	Count    int            `json:"count"`
	Current  int            `json:"current"`
	Next     int            `json:"next"`
	Previous int            `json:"previous"`
	Last     int            `json:"last"`
	Results  []EventDetails `json:"results"`
}

type EventDetails struct {
	ID              string                 `json:"id"`
	Shell           string                 `json:"shell"`
	Line            string                 `json:"line"`
	Success         bool                   `json:"success"`
	Result          string                 `json:"result"`
	Status          map[string]interface{} `json:"status"`
	ResponseDelay   float64                `json:"response_delay"`
	ElapsedTime     float64                `json:"elapsed_time"`
	AddedAt         time.Time              `json:"added_at"`
	Server          string                 `json:"server"`
	ServerName      string                 `json:"server_name"`
	RequestedBy     string                 `json:"requested_by"`
	RequestedByName string                 `json:"requested_by_name"`
}
