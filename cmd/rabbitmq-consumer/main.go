package main

import (
	"context"
	"database/sql"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/vfunin/rabbitmq-consumer/internal/config"
	"github.com/wagslane/go-rabbitmq"
)

type contextKey int

const (
	contextKeyID contextKey = iota
)

func main() {
	var cPath string

	rand.Seed(time.Now().UnixNano())

	logger, conf, err := getConfigAndLog(cPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Start application...")

	driver := "postgres"

	if !strings.Contains(conf.DatabaseDSN, "postgres") {
		driver = "mysql"
	}

	db, err := sql.Open(driver, conf.DatabaseDSN)
	if err != nil {
		log.Fatalf("Connection DB error: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to reach database: %v", err)
	}

	db.SetMaxOpenConns(conf.DBMaxOpenConns)
	db.SetMaxIdleConns(conf.DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(conf.DBConnMaxLifetime) * time.Minute)

	defer db.Close()

	consumer, err := rabbitmq.NewConsumer(
		conf.RabbitDSN,
		rabbitmq.Config{ //nolint:exhaustivestruct
			Heartbeat: 10 * time.Second, //nolint:gomnd
			Locale:    "ru_RU",
		},
		rabbitmq.WithConsumerOptionsLogger(logger),
		rabbitmq.WithConsumerOptionsReconnectInterval(time.Duration(conf.RabbitReconnectInterval)*time.Second),
	)

	if err != nil {
		logger.Fatalf("Connection error: %v", err)
	}

	defer consumer.Close()

	err = consumer.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			sqlString := "INSERT INTO " + conf.TableName + " (message, correlation_id) VALUES "
			if driver == "mysql" {
				_, err = db.Exec(sqlString+"(?, ?)", string(d.Body), d.CorrelationId)
			} else {
				_, err = db.Exec(sqlString+"($1, $2)", string(d.Body), d.CorrelationId)
			}

			if err != nil {
				log.Fatalf("An error occurred while executing query: %v", err)
			}

			return rabbitmq.Ack
		},
		conf.Queue,
		[]string{"routing_key", "routing_key_2"},
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsConsumerName(conf.ConsumerName),
		rabbitmq.WithConsumeOptionsConcurrency(conf.RabbitGoroutinesCnt),
	)

	if err != nil {
		logger.Fatalf("Consuming error: %v", err)
	}

	ctx, cancel := getContext(logger)

	listenChannels(ctx, cancel)
}

func getConfigAndLog(cPath string) (*log.Entry, config.Config, error) {
	log.SetFormatter(&log.TextFormatter{ //nolint:exhaustivestruct
		DisableColors: true,
		FullTimestamp: true,
	})

	logger := log.WithFields(log.Fields{
		"service": "consumer",
	})

	flag.StringVar(&cPath, "config", "", "path to config.yaml")
	flag.Parse()

	conf, err := config.GetConfig(cPath)

	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}

	if conf.LogFormat != "text" {
		log.SetFormatter(&log.JSONFormatter{}) //nolint:exhaustivestruct
	}

	if conf.ConsumerName == "go-consumer" {
		conf.ConsumerName += "-" + strconv.Itoa(rand.Intn(9999-1000)+1000) //nolint:gomnd,gosec
	}

	return logger, conf, err
}

func getContext(logger *log.Entry) (ctx context.Context, cancel context.CancelFunc) {
	cwv := context.WithValue(context.Background(), contextKeyID, logger)

	ctx, cancel = context.WithCancel(cwv)

	return
}

func listenChannels(ctx context.Context, cancel context.CancelFunc) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT)

	incCh := make(chan os.Signal, 1)
	signal.Notify(incCh, syscall.SIGUSR1)

	logger := ctx.Value(contextKeyID).(*log.Entry)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Consuming done")

			return
		case <-stopCh:
			logger.Info("Start graceful shutdown")
			cancel()
		case <-incCh:
			logger.Info("-")
		}
	}
}
