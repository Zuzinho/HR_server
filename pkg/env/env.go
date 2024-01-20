package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// LinksConfig stored links to
// Agify - age forecaster
// Genderize - gender forecaster
// Nationalize - nation forecaster
type LinksConfig struct {
	AgifyLink       string
	GenderizeLink   string
	NationalizeLink string
}

func init() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	log.Println("loaded .env file")
}

// MustDBConnString returns db connection string from environment
// By key DB_CONNECTION_STRING
func MustDBConnString() string {
	val, exist := os.LookupEnv("DB_CONNECTION_STRING")
	if !exist {
		log.Fatal("no DB connection string")
	}

	log.Println("looked up DB connection string")
	return val
}

// MustLinksConfig returns LinksConfig from environment
// By keys AGIFY_LINK, GENDERIZE_LINK, NATIONALIZE_LINK
func MustLinksConfig() *LinksConfig {
	agify, exist := os.LookupEnv("AGIFY_LINK")
	if !exist {
		log.Fatal("no agify link")
	}

	log.Println("looked up agify link")

	genderize, exist := os.LookupEnv("GENDERIZE_LINK")
	if !exist {
		log.Fatal("no genderize link")
	}

	log.Println("looked up genderize link")

	nationalize, exist := os.LookupEnv("NATIONALIZE_LINK")
	if !exist {
		log.Fatal("no nationalize link")
	}

	log.Println("looked up nationalize link")

	return &LinksConfig{
		AgifyLink:       agify,
		GenderizeLink:   genderize,
		NationalizeLink: nationalize,
	}
}

// MustPort returns port from environment
// By key PORT
func MustPort() string {
	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("no port")
	}

	log.Println("looked up port")

	return port
}

// MustMaxOpenConns returns max open connections from environment
// By key MAX_OPEN_CONNS
func MustMaxOpenConns() int {
	str, exist := os.LookupEnv("MAX_OPEN_CONNS")
	if !exist {
		log.Fatal("no max open conns")
	}

	log.Println("looked up max open conns")

	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	return val
}
