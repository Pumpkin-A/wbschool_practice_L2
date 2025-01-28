package config

type Config struct {
	DbFilename string
	Port       string
}

func NewDefaultConfig() *Config {
	return &Config{DbFilename: "db.txt", Port: ":8080"}
}

func NewTestConfig() *Config {
	return &Config{DbFilename: "test_db.txt", Port: ":8081"}
}
