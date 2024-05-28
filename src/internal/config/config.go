package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Model           `yaml:"model"`
	ReportStorage   `yaml:"report_storage"`
	DocumentStorage `yaml:"document_storage"`
	Database        `yaml:"database"`
	HTTPServer      `yaml:"http_server"`
	Logger          `yaml:"logger"`
}

type Model struct {
	Route string `yaml:"route"`
}

type ReportStorage struct {
	ReportCreatorPath string `yaml:"report_creator_path"`
	ReportPath        string `yaml:"report_storage"`
	ReportExt         string `yaml:"report_ext"`
}

type DocumentStorage struct {
	DocumentPath string `yaml:"document_path"`
	DocumentExt  string `yaml:"document_ext"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Database string `yaml:"database" env-required:"true"`
	Port     string `yaml:"port"`
}

type HTTPServer struct {
	Addr         string        `yaml:"addr" env-default:"localhost:8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"40s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"40s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"40s"`
}

type Logger struct {
	LogLevel        string `yaml:"log_level"`
	OutputFilePath  string `yaml:"output_filepath"`
	UseFile         bool   `yaml:"use_file" env-default:"false"`
	LogFormat       string `yaml:"log_format"`
	TimestampFormat string `yaml:"timestamp_format"`
	OutputFormat    string `yaml:"output_format"`
}

func (db *Database) GetGormConnectStr() string {
	return fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s", db.Host, db.User, db.Password, db.Database, db.Port)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
