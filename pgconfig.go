package easydb

import "strconv"

type PGConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Params   map[string]string
}

// NewPGConfig создаёт базовый конфиг
func NewPGConfig() PGConfig {
	return PGConfig{
		Host:   "localhost",
		Port:   "5432",
		Params: make(map[string]string),
	}
}

// Обёртки

func (cfg PGConfig) WithHost(host string) PGConfig {
	cfg.Host = host
	return cfg
}

func (cfg PGConfig) WithPort(port string) PGConfig {
	cfg.Port = port
	return cfg
}

func (cfg PGConfig) WithUser(user string) PGConfig {
	cfg.Username = user
	return cfg
}

func (cfg PGConfig) WithPassword(pw string) PGConfig {
	cfg.Password = pw
	return cfg
}

func (cfg PGConfig) WithDatabase(name string) PGConfig {
	cfg.DBName = name
	return cfg
}

func (cfg PGConfig) WithSSLMode(mode string) PGConfig {
	cfg.Params["sslmode"] = mode
	return cfg
}

func (cfg PGConfig) WithMaxConns(n int) PGConfig {
	cfg.Params["pool_max_conns"] = strconv.Itoa(n)
	return cfg
}

func (cfg PGConfig) WithAppName(name string) PGConfig {
	cfg.Params["application_name"] = name
	return cfg
}

func (cfg PGConfig) WithSchema(schema string) PGConfig {
	cfg.Params["search_path"] = schema
	return cfg
}
