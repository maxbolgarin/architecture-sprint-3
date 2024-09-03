package main

import "time"

type TelemetryModel struct {
	Value float64   `db:"value"`
	Time  time.Time `db:"time"`
}

type TelemetrySchema struct {
	MetricName   string         `json:"metric_name"`
	MetricValues []MetricValues `json:"values"`
}

type MetricValues struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type TelemetryMessage struct {
	DeviceID string    `json:"device_id"`
	MetricID string    `json:"metric_id"`
	Value    float64   `json:"value"`
	Time     time.Time `json:"time"`
}
