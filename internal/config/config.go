package config

import (
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort string

	HttpRateLimitRequest int
	HttpRateLimitTime    time.Duration

	JwtSecretKey string
	JwtTTL       time.Duration

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

func panicIfNeeded(err error, i int) {
	if err != nil {
		log.Println(err.Error() + " : LINE", i)
	}
}

func load() Config {
	AppPort := os.Getenv("SERVER_URL")

	HttpRateLimitRequest, err := strconv.Atoi(os.Getenv("HttpRateLimitRequest"))
	panicIfNeeded(err, 50)

	HttpRateLimitTime,err := time.ParseDuration(os.Getenv("HttpRateLimitTime"))
	panicIfNeeded(err, 53)

	JwtSecretKey := os.Getenv("JwtSecretKey")
	JwtTTL, err := time.ParseDuration(os.Getenv("JwtTTL"))
	panicIfNeeded(err, 57)

	PaginationLimit, err := strconv.Atoi(os.Getenv("PaginationLimit"))
	panicIfNeeded(err, 60)

	MysqlUser := os.Getenv("MysqlUser")
	MysqlPassword := os.Getenv("MysqlPassword")
	MysqlHost := os.Getenv("MysqlHost")
	MysqlPort,err := strconv.Atoi(os.Getenv("MysqlPort"))
	panicIfNeeded(err, 66)

	MysqlDatabase := os.Getenv("MysqlDatabase")
	MysqlMaxIdleConns,err := strconv.Atoi(os.Getenv("MysqlMaxIdleConns"))
	panicIfNeeded(err, 70)

	MysqlMaxOpenConns,err := strconv.Atoi(os.Getenv("MysqlMaxOpenConns"))
	panicIfNeeded(err, 73)

	MysqlConnMaxLifetime,err := time.ParseDuration(os.Getenv("MysqlConnMaxLifetime"))
	panicIfNeeded(err, 76)

	RedisPassword := os.Getenv("RedisPassword")
	RedisHost := os.Getenv("RedisHost")
	RedisPort,err := strconv.Atoi(os.Getenv("RedisPort"))
	panicIfNeeded(err, 81)

	RedisDatabase, err := strconv.Atoi(os.Getenv("RedisDatabase"))
	panicIfNeeded(err, 84)

	RedisPoolSize, err := strconv.Atoi(os.Getenv("RedisPoolSize"))
	panicIfNeeded(err, 87)

	RedisTTL,err := time.ParseDuration(os.Getenv("RedisTTL"))
	panicIfNeeded(err, 90)

	return Config{
		AppPort:              AppPort,
		HttpRateLimitRequest: HttpRateLimitRequest,
		HttpRateLimitTime:    HttpRateLimitTime,
		JwtSecretKey:         JwtSecretKey,
		JwtTTL:               JwtTTL,
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