package configs

import "time"

type Pagelines struct {
	PageSizeKey     string `yaml:"PageSizeKey"`
	DefaultPage     int    `yaml:"DefaultPage"`
	DefaultPageSize int    `yaml:"DefaultPageSize"`
	PageKey         string `yaml:"PageKey"`
}

type Server struct {
	RunMode               string        `yaml:"RunMode"`
	Address               string        `yaml:"Address"`
	ReadTimeout           time.Duration `yaml:"ReadTimeout"`
	WriteTimeout          time.Duration `yaml:"WriteTimeout"`
	DefaultContextTimeout time.Duration `yaml:"DefaultContextTimeout"`
}

type App struct {
	Version string `yaml:"Version"`
	Name    string `yaml:"Name"`
}

type Log struct {
	LogFileExt    string `yaml:"LogFileExt"`
	MaxAge        int    `yaml:"MaxAge"`
	MaxBackups    int    `yaml:"MaxBackups"`
	Level         string `yaml:"Level"`
	HighLevelFile string `yaml:"HighLevelFile"`
	LowLevelFile  string `yaml:"LowLevelFile"`
	MaxSize       int    `yaml:"MaxSize"`
	Compress      bool   `yaml:"Compress"`
	LogSavePath   string `yaml:"LogSavePath"`
}

type Email struct {
	Host     string   `yaml:"Host"`
	Port     int      `yaml:"Port"`
	UserName string   `yaml:"UserName"`
	Password string   `yaml:"Password"`
	IsSSL    bool     `yaml:"IsSSL"`
	From     string   `yaml:"From"`
	To       []string `yaml:"To"`
}

type Postgres struct {
	DBDriver string `yaml:"DBDriver"`
	Address  string `yaml:"Address"`
	DBName   string `yaml:"DBName"`
	Sslmode  string `yaml:"Sslmode"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
}

type Token struct {
	Key               string        `yaml:"key"`
	Duration          time.Duration `yaml:"Duration"`
	AuthorizationKey  string        `yaml:"AuthorizationKey"`
	AuthorizationType string        `yaml:"AuthorizationType"`
}
