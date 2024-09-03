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
	_, err := Exec(ctx, db, createCommandsTableQuery)
	if err != nil {
		return err
	}

	mux := mux.NewRouter()
	out := &app{
		db: db,
	}

	mux.HandleFunc("/api/v1/commands/{device_id}", out.getForDevice)
	mux.HandleFunc("/api/v1/commands/history/{device_id}", out.getHistorical)

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

func (s *app) getForDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceID := params["device_id"]

	commands, err := QueryAll[CommandModel](r.Context(), s.db, getForDeviceQuery, deviceID)
	switch {
	case errors.Is(err, ErrEmptyResult):
		httpErr(w, http.StatusNotFound, err)
		return

	case err != nil:
		httpErr(w, http.StatusInternalServerError, err)
		return
	}

	for i := range commands {
		_, err = Exec(r.Context(), s.db, updateSendQuery, deviceID, commands[i].CommandID)
		if err != nil {
			httpErr(w, http.StatusInternalServerError, err)
			return
		}
	}

	out := make([]CommandSchema, 0, len(commands))
	for _, t := range commands {
		out = append(out, CommandSchema{
			CommandID:   t.CommandID,
			CreateTime:  t.CreateTime,
			Code:        t.Code,
			CommandType: t.CommandType,
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

func (s *app) getHistorical(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceID := params["device_id"]
	from := r.URL.Query().Get("from")

	if from == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fromInt, err := strconv.Atoi(from)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fromTime := time.Unix(int64(fromInt), 0)

	commands, err := QueryAll[CommandModel](r.Context(), s.db, getHistoryQuery, deviceID, fromTime)
	switch {
	case errors.Is(err, ErrEmptyResult):
		httpErr(w, http.StatusNotFound, err)
		return

	case err != nil:
		httpErr(w, http.StatusInternalServerError, err)
		return
	}

	out := make([]CommandHistorySchema, 0, len(commands))
	for _, t := range commands {
		out = append(out, CommandHistorySchema{
			CommandID:     t.CommandID,
			CreateTime:    t.CreateTime,
			Code:          t.Code,
			CommandType:   t.CommandType,
			SendTime:      t.SendTime,
			CommandTypeID: t.CommandTypeID,
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
	getForDeviceQuery        = "SELECT (command_id, create_time, command_type, code) FROM commands WHERE device_id = $1 AND send_time IS NULL"
	getHistoryQuery          = "SELECT (command_id, command_type_id, create_time, send_time, command_type, code) FROM commands WHERE device_id = $1 AND create_time > $2"
	updateSendQuery          = "UPDATE commands SET send_time = now() WHERE device_id = $1 AND command_id = $2"
	createCommandsTableQuery = `
	CREATE TABLE IF NOT EXISTS commands (
		command_id TEXT NOT NULL,
		device_id TEXT NOT NULL,
		command_type_id TEXT NOT NULL,
		create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		send_time TIMESTAMPTZ,
		command_type TEXT NOT NULL,
		code TEXT NOT NULL,
		PRIMARY KEY (command_id)
	)
	`
)
