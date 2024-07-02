package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/hibiken/asynq"
	"github.com/tiago123456789/poc-collect-log-golang/model"
)

const redisAddr = "127.0.0.1:6379"

func main() {

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 3,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("log", func(c context.Context, task *asynq.Task) error {
		fmt.Println("starting to process message")
		var logData model.Log
		err := json.Unmarshal(task.Payload(), &logData)
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Error: %v", err),
			)
			return err
		}

		logData.Level = strings.ToLower(logData.Level)
		logToByte, err := json.Marshal(logData)
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Error: %v", err),
			)
			return err
		}

		req := esapi.IndexRequest{
			Index: "logs",
			// DocumentID: strconv.Itoa(1),
			Body:    bytes.NewReader(logToByte),
			Refresh: "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document ID", res.Status())
		}
		return nil
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
