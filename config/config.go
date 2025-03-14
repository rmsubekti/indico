package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB struct {
		USER        string
		PASSWORD    string
		NAME        string
		HOSTNAME    string
		PORT        string
		SSLMODE     string
		AUTH_METHOD string
	}

	APP struct {
		LOG_LEVEL    string
		PORT         string
		HOSTNAME     string
		JWT_SECRET   string
		GATEWAY_HOST string
	}
)

func init() {
	godotenv.Load()

	// db config
	DB.USER = get("DB_USER", "indico")
	DB.PASSWORD = get("DB_PASSWORD", "indico-00BC7ddAa54cpB")
	DB.NAME = get("DB_NAME", "indico")
	DB.HOSTNAME = get("DB_HOSTNAME", "localhost")
	DB.PORT = get("DB_PORT", "5432")
	DB.SSLMODE = get("DB_SSLMODE", "disable")
	DB.AUTH_METHOD = get("DB_AUTH_METHOD", "trust")

	//app config
	APP.LOG_LEVEL = get("APP_LOG_LEVEL", "debug")
	APP.PORT = get("APP_PORT", "2011")
	APP.HOSTNAME = get("APP_HOSTNAME", "localhost")
	APP.JWT_SECRET = get("APP_JWT_SECRET", "'09)77)45e;sd3e9;92d__=12'")
	APP.GATEWAY_HOST = get("APP_GATEWAY_HOST", "localhost:3300")
}

func PG_DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		DB.HOSTNAME, DB.USER, DB.PASSWORD, DB.NAME, DB.PORT, DB.SSLMODE)
}

func PG_URI() (uri string) {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		DB.USER, DB.PASSWORD, DB.HOSTNAME, DB.PORT, DB.NAME)
}

func get(key, defaultVal string) string {
	v := os.Getenv(key)
	if len(v) > 0 {
		return v
	}
	return defaultVal
}
