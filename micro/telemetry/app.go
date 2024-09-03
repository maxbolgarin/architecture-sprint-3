package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func StartApp(ctx context.Context, db *DB) error {
	_, err := Exec(ctx, db, createTelemetryTableQuery)
	if err != nil {
		return err
	}

	mux := mux.NewRouter()
	out := &app{
		db: db,
	}

	mux.HandleFunc("/api/v1/telemetry/{device_id}/{metric_id}", out.get)

	srv := http.Server{
		Addr:    os.Getenv("SERVER_ADDRESS"),
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("starting server at " + srv.Addr)

	return nil
}

type app struct {
	db *DB
}

func (s *app) get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceID := params["device_id"]
	metricID := params["metric_id"]
	from := r.URL.Query().Get("from")

	if from == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fromInt, err := strconv.Atoi(from)
	if err != nil {
		httpErr(w, http.StatusBadRequest, err)
		return
	}

	fromTime := time.Unix(int64(fromInt), 0)

	telemetry, err := QueryAll[TelemetryModel](r.Context(), s.db, getTelemetryQuery, deviceID, metricID, fromTime)
	switch {
	case errors.Is(err, ErrEmptyResult):
		httpErr(w, http.StatusNotFound, err)
		return
	case err != nil:
		httpErr(w, http.StatusInternalServerError, err)
		return
	}

	out := TelemetrySchema{
		MetricName:   metricID, // TODO: request to devices types to get name
		MetricValues: make([]MetricValues, 0, len(telemetry)),
	}

	for _, t := range telemetry {
		out.MetricValues = append(out.MetricValues, MetricValues{
			Time:  t.Time.Unix(),
			Value: t.Value,
		})
	}

	outBytes, err := json.Marshal(out)
	if err != nil {
		httpErr(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(outBytes)
}

func httpErr(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	fmt.Fprintln(w, err.Error())
	log.Println("http error:", err.Error())
}

const (
	getTelemetryQuery         = "SELECT (value, create_time) FROM telemetry WHERE device_id = $1 AND metric_id = $2 AND create_time >= $3"
	createTelemetryTableQuery = `
	CREATE TABLE IF NOT EXISTS telemetry (
		device_id TEXT NOT NULL,
		metric_id TEXT NOT NULL,
		create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		value INT NOT NULL,
		PRIMARY KEY (device_id, metric_id, create_time)
	)
	`
)
