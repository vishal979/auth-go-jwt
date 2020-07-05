package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	// mysql driver import
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

// Server driver struct
type Server struct {
	DB          *sql.DB
	RedisClient *redis.Client
	Router      *mux.Router
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (server *Server) waitToConnect() error {
	var err error
	log.Error("trying connecting to database")
	for i := 0; i < 60; i++ {
		log.Info("trying ", i+1)
		err = server.DB.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}

//Initialize intializes database connections
func (server *Server) Initialize() {
	var err error

	DbUser := getEnv("DB_USER", "root")
	DbPassword := getEnv("DB_PASSWORD", "vishal1132")
	DbHost := getEnv("MYSQL_URL", "localhost")
	DbPort := getEnv("DB_PORt", "3306")
	DbName := getEnv("DB_NAME", "testdb")
	//mysql connection and migration
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = sql.Open(getEnv("DB_DRIVER", "mysql"), DBURL)
	if err != nil {
		log.Fatal("error connecting to mysql database", err)
	}
	if err = server.waitToConnect(); err != nil {
		log.Fatal("error pinging mysql ", err)
	}
	log.Info("Mysql DB Connection Successful")
	server.migrateTables()
	server.RedisClient = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
	_, err = server.RedisClient.Ping().Result()
	if err != nil {
		log.Fatal("Error connecting to redis server", err)
	}
	log.Info("Connection to redis successful")
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

// Run runs server
func (server *Server) Run(addr string) {
	log.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) migrateTables() {
	log.Info("Migrating Schema")
	server.DB.Exec("use testdb")
	server.DB.Exec("CREATE TABLE IF NOT EXISTS Users(`id` bigint NOT NULL AUTO_INCREMENT,`username` varchar(100) NOT NULL UNIQUE,`password` varchar(400) NOT NULL,PRIMARY KEY (`id`),INDEX(`username`))  ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;")
}
