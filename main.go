package main

import (
	"flag"
	"github.com/mqiqe/prometheus-m3db-sarama/pkg/saramaservice"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	brokers  = "localhost:9092"
	group    = "metrics"
	topics   = "metrics"
	storeUrl = "http://127.0.0.1:7201/api/v1/prom/remote/write"
	version  = "2.1.1"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	flag.StringVar(&brokers, "brokers", "127.0.0.1:9092", "Kafka Address")
	flag.StringVar(&group, "group", "metrics", "Kafka Group")
	flag.StringVar(&topics, "topics", "metrics", "Kafka Topics")
	flag.StringVar(&storeUrl, "store-url", "http://127.0.0.1:7201/api/v1/prom/remote/write", "M3db Store Url")
	flag.Parse()

	if value := os.Getenv("KAFKA_BROKERS"); value != "" {
		brokers = value
	}
	if value := os.Getenv("KAFKA_GROUP"); value != "" {
		group = value
	}
	if value := os.Getenv("KAFKA_TOPICS"); value != "" {
		topics = value
	}
	if value := os.Getenv("STORE_RUL"); value != "" {
		storeUrl = value
	}
}
func main() {
	// 接收信息
	log.Infof("start : %v", time.Now())
	ks := saramaservice.NewSaramaService(brokers, group, topics, version, storeUrl)
	ks.Run()
	log.Infof("stop : %v", time.Now())
}
