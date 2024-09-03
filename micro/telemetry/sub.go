package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	TelemetryTopic = "telemetry"

	InsertTelemetryQuery = "INSERT INTO telemetry (device_id, metric_id, value, created_time) VALUES ($1, $2, $3, $4)"
)

func StartKafka(ctx context.Context, db *DB) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_ADDRESS")},
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: os.Getenv("KAFKA_GROUP"),
	})

	go func() {
		for {
			kafkaMsg, err := r.ReadMessage(ctx)
			if err != nil {
				log.Println("cannot read message from kafka: " + err.Error())
				continue
			}

			var msg TelemetryMessage
			if err := json.Unmarshal(kafkaMsg.Value, &msg); err != nil {
				log.Println("cannot unmarshal from kafka: " + err.Error())
				continue
			}

			_, err = Exec(ctx, db, InsertTelemetryQuery, msg.DeviceID, msg.MetricID, msg.Value, msg.Time)
			if err != nil {
				log.Println("cannot insert msg from kafka to DB: " + err.Error())
				continue
			}
			fmt.Println("inserted a new telemetry for device " + msg.DeviceID)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}

	}()
	go func() {
		<-ctx.Done()
		if err := r.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()

	log.Println("starting kafka sub at " + os.Getenv("KAFKA_ADDRESS"))

	return nil
}
