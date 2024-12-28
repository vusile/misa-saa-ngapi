package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	router   http.Handler
	db       *sql.DB
	gorm     *gorm.DB
	config   Config
	esClient *elasticsearch.Client
}

var AppInstance *App

func New(config Config) *App {

	params := make(map[string]string)
	params["parseTime"] = "true"

	cfg := mysql.Config{
		User:   config.DBUsername,
		Passwd: config.DBPassword,
		Net:    "tcp",
		Addr:   config.DBPort,
		DBName: config.DBName,
		Params: params,
	}

	//todo: handle error
	db, _ := sql.Open("mysql", cfg.FormatDSN())

	//todo: handle error
	gormDB, _ := gorm.Open(gormmysql.New(gormmysql.Config{
		Conn: db,
	}), &gorm.Config{})

	client, err := ConnectToESClient(config)

	if err != nil {
		fmt.Println("failed to connect to elastic search: %w", err)
	}

	AppInstance := &App{
		db:       db,
		gorm:     gormDB,
		config:   config,
		esClient: client,
	}

	AppInstance.loadRoutes()

	return AppInstance
}

func (a *App) Start(ctx context.Context) error {
	err := a.db.Ping()

	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}

	migrate(a)
	CreateESIndex(a.esClient)
	IndexParokia(a)
	SearchAsYouType(a.esClient)

	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("failed to close mysql", err)
		}
	}()

	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	fmt.Println("Starting Server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}

func GetApplicationInstance() *App {
	return AppInstance
}
