package configs

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port         int
	DatabaseURL  string
	KafkaBroker  string
	KafkaTopics  []string
	KafkaGroupID string
}

func LoadConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalln("PORT environment variable is not set")
	}
	cfg := &Config{
		Port:         port,
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		KafkaBroker:  os.Getenv("KAFKA_BROKER"),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
	}

	topicsStr := os.Getenv("KAFKA_TOPICS")
	cfg.KafkaTopics = strings.Split(topicsStr, ",")
	// Trim spaces from each topic name
	for i, topic := range cfg.KafkaTopics {
		cfg.KafkaTopics[i] = strings.TrimSpace(topic)
	}

	if cfg.DatabaseURL == "" {
		log.Fatalln("DATABASE_URL environment variable is not set")
	}

	if cfg.KafkaBroker == "" {
		log.Fatalln("KAFKA_BROKER environment variable is not set")
	}

	if cfg.KafkaGroupID == "" {
		log.Fatalln("KAFKA_GROUP_ID environment variable is not set")
	}

	return cfg
}
