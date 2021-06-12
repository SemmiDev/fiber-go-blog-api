package config

import (
	"github.com/SemmiDev/fiber-go-blog/internal/logger"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort string

	HttpRateLimitRequest int
	HttpRateLimitTime    time.Duration

	PaginationLimit int

	MysqlUser            string
	MysqlPassword        string
	MysqlHost            string
	MysqlPort            int
	MysqlDatabase        string
	MysqlMaxIdleConns    int
	MysqlMaxOpenConns    int
	MysqlConnMaxLifetime time.Duration

	RedisPassword string
	RedisHost     string
	RedisPort     int
	RedisDatabase int
	RedisPoolSize int
	RedisTTL      time.Duration
}

func load() Config {
	AppPort := os.Getenv("SERVER_URL")

	HttpRateLimitRequest, err := strconv.Atoi(os.Getenv("HttpRateLimitRequest"))
	logger.Log().Err(err)

	HttpRateLimitTime,err := time.ParseDuration(os.Getenv("HttpRateLimitTime"))
	logger.Log().Err(err)

	PaginationLimit, err := strconv.Atoi(os.Getenv("PaginationLimit"))
	logger.Log().Err(err)

	MysqlUser := os.Getenv("MysqlUser")
	MysqlPassword := os.Getenv("MysqlPassword")
	MysqlHost := os.Getenv("MysqlHost")
	MysqlPort,err := strconv.Atoi(os.Getenv("MysqlPort"))
	logger.Log().Err(err)

	MysqlDatabase := os.Getenv("MysqlDatabase")
	MysqlMaxIdleConns,err := strconv.Atoi(os.Getenv("MysqlMaxIdleConns"))
	logger.Log().Err(err)

	MysqlMaxOpenConns,err := strconv.Atoi(os.Getenv("MysqlMaxOpenConns"))
	logger.Log().Err(err)

	MysqlConnMaxLifetime,err := time.ParseDuration(os.Getenv("MysqlConnMaxLifetime"))
	logger.Log().Err(err)

	RedisPassword := os.Getenv("RedisPassword")
	RedisHost := os.Getenv("RedisHost")
	RedisPort,err := strconv.Atoi(os.Getenv("RedisPort"))
	logger.Log().Err(err)

	RedisDatabase, err := strconv.Atoi(os.Getenv("RedisDatabase"))
	logger.Log().Err(err)

	RedisPoolSize, err := strconv.Atoi(os.Getenv("RedisPoolSize"))
	logger.Log().Err(err)

	RedisTTL,err := time.ParseDuration(os.Getenv("RedisTTL"))
	logger.Log().Err(err)

	return Config{
		AppPort:              AppPort,
		HttpRateLimitRequest: HttpRateLimitRequest,
		HttpRateLimitTime:    HttpRateLimitTime,
		PaginationLimit:      PaginationLimit,
		MysqlUser:            MysqlUser,
		MysqlPassword:        MysqlPassword,
		MysqlHost:            MysqlHost,
		MysqlPort:            MysqlPort,
		MysqlDatabase:        MysqlDatabase,
		MysqlMaxIdleConns:    MysqlMaxIdleConns,
		MysqlMaxOpenConns:    MysqlMaxOpenConns,
		MysqlConnMaxLifetime: MysqlConnMaxLifetime,
		RedisPassword:        RedisPassword,
		RedisHost:            RedisHost,
		RedisPort:            RedisPort,
		RedisDatabase:        RedisDatabase,
		RedisPoolSize:        RedisPoolSize,
		RedisTTL:             RedisTTL,
	}
}

var config = load()

func Cfg() *Config { return &config }