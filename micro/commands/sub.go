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
	insertCommandQuery = "INSERT INTO commands (command_id, command_type_id, create_time, send_time, command_type, code) VALUES ($1, $2, $3, $4, $5, $6)"
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

			var msg CommandModel
			if err := json.Unmarshal(kafkaMsg.Value, &msg); err != nil {
				log.Println("cannot unmarshal from kafka: " + err.Error())
				continue
			}

			_, err = Exec(ctx, db, insertCommandQuery, msg.CommandID, msg.CommandTypeID, msg.CreateTime, msg.SendTime, msg.CommandType, msg.Code)
			if err != nil {
				log.Println("cannot insert msg from kafka to DB: " + err.Error())
				continue
			}
			fmt.Println("inserted a new telemetry for device " + msg.CommandID)

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
