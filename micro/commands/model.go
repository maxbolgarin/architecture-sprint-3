package main

import "time"

type CommandModel struct {
	CommandID     string    `db:"command_id"`
	CommandTypeID string    `db:"command_type_id"`
	CreateTime    time.Time `db:"create_time"`
	SendTime      time.Time `db:"send_time"`
	Code          string    `db:"code"`
	CommandType   string    `db:"command_type"`
}

type CommandSchema struct {
	CommandID   string    `json:"command_id"`
	CreateTime  time.Time `json:"create_time"`
	Code        string    `json:"code"`
	CommandType string    `json:"command_type"`
}

type CommandHistorySchema struct {
	CommandID     string    `json:"command_id"`
	CommandTypeID string    `json:"command_type_id"`
	CreateTime    time.Time `json:"create_time"`
	SendTime      time.Time `json:"send_time"`
	Code          string    `json:"code"`
	CommandType   string    `json:"command_type"`
}
