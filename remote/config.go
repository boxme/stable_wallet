package remote

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port     int      `json:"port"`
	Env      string   `json: "env"`
	Database DbConfig `json: "database:"`
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name`
}

func LoadConfig(isProd bool) Config {
	if !isProd {
		fmt.Println("Successfully loaded dev config")
		return devConfig()
	}

	f, err := os.Open(".config")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(f)
	var c Config
	err = decoder.Decode(&c)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully loaded prod config")
	return c
}

func devConfig() Config {
	dbConfig := getDevDbConfig()
	return Config{
		Port:     3000,
		Env:      "dev",
		Database: dbConfig,
	}
}

func getDevDbConfig() DbConfig {
	return DbConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "dev",
		Password: "honeybbee8988",
		DbName:   "stable_wallet",
	}
}

func (c *DbConfig) GetDbConnectionString() string {
	// databaseUrl := "postgres://postgres:mypassword@localhost:5432/postgres"
	return fmt.Sprintf(
		"%s://%s:%s@localhost:%d/%s",
		c.Host,
		c.User,
		c.Password,
		c.Port,
		c.DbName)
}
