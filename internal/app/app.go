package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tests2/internal/cache"
	"tests2/internal/config"
	"tests2/internal/http"
	"tests2/internal/minio"
	"tests2/internal/models"
	"tests2/internal/repository"
	"tests2/internal/repository/client"
	"tests2/internal/server"
	"tests2/internal/service"
	"time"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(io.MultiWriter(writers...))

	db, err := sql.Open(cfg.DriverName, cfg.DBConnStr)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	pgClient := client.NewPostgresClient(db)

	repos := repository.NewRepositories(&pgClient)

	min, err := minio.NewMinIOClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := min.InitDocBucket(context.Background()); err != nil {
		log.Fatal(err)
	}

	services := service.NewService(repos, cfg, min)

	c, err := cache.NewCache(1024)
	if err != nil {
		logrus.Fatal("Cannot init cache ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer cancel()
	orders, err := services.Order.GetOrderList(ctx)
	if err != nil {
		logrus.Fatal("Cannot upload cache", err)
	}
	err = c.UploadCache(orders)
	if err != nil {
		logrus.Fatal("cannot upload cache")
	}

	srv := server.NewServer(cfg)
	http.NewHandlers(cfg, services, c).Init(srv.App())

	go func() {
		err := srv.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	//sc, _ := stan.Connect("test-cluster", "sub-1", stan.NatsURL("nats://0.0.0.0:4223"))
	//defer sc.Close()
	//
	//go func() {
	//	sc.Subscribe("orders", func(m *stan.Msg) {
	//		order := &models.Order{}
	//
	//		err := json.Unmarshal(m.Data, &order)
	//		if err != nil {
	//			fmt.Println("Cannot Marshal import")
	//		}
	//
	//		err = validateOrder(order)
	//		if err != nil {
	//			logrus.Error("Order is not valid")
	//		} else {
	//			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	//			defer cancel()
	//			err = services.Order.CreateOrder(ctx, order)
	//			if err != nil {
	//				logrus.Error(err)
	//			}
	//			ok := c.Add(order.OrderUID, order)
	//			if ok {
	//				fmt.Println("eviction occurred")
	//			} else {
	//				fmt.Println("eviction didn't occur")
	//			}
	//		}
	//	})
	//}()
	//
	//Block()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := srv.Stop(); err != nil {
		logrus.Fatal("Server forced to shut down", err)
	}
}

func Convert(input interface{}) ([]byte, error) {
	var order []byte
	err := mapstructure.Decode(input, &order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func Block() {
	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}

var Validator = validator.New()

func validateOrder(ord *models.Order) error {
	err := Validator.Struct(ord)
	if err != nil {
		return errors.New("cannot validate order")
	}
	return nil
}
