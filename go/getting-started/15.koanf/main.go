package main

import (
	"fmt"
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

// Config represents the structure of the configuration file
type Config struct {
	AppName   string `koanf:"app_name"`
	Port      int    `koanf:"port"`
	DebugMode bool   `koanf:"debug_mode"`
	Database  struct {
		Host     string `koanf:"host"`
		Port     int    `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		DBName   string `koanf:"db_name"`
	} `koanf:"database"`
}

func main() {
	k := koanf.New(".")

	// Load the main configuration
	err := k.Load(file.Provider("config.json"), json.Parser())
	if err != nil {
		log.Fatalf("Error loading main configuration: %v", err)
	}

	// Load the environment-specific configuration
	err = k.Load(file.Provider("config.local.json"), json.Parser())
	if err != nil {
		log.Printf("Warning: could not load environment configuration: %v", err)
	}

	// Unmarshal merged configuration into Config struct
	var finalConfig Config
	err = k.Unmarshal("", &finalConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling configuration: %v", err)
	}

	fmt.Printf("App Name: %s\n", finalConfig.AppName)
	fmt.Printf("Port: %d\n", finalConfig.Port)
	fmt.Printf("Debug Mode: %v\n", finalConfig.DebugMode)
	fmt.Printf("Database Host: %s\n", finalConfig.Database.Host)
	fmt.Printf("Database Port: %d\n", finalConfig.Database.Port)
	fmt.Printf("Database User: %s\n", finalConfig.Database.User)
	fmt.Printf("Database Name: %s\n", finalConfig.Database.DBName)
}
