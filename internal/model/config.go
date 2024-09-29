package model

type ServiceConfig struct {
	Host     string         `yaml:"host"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	DBName       string `yaml:"db-name"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	PoolMaxConns int    `yaml:"pool-max-conns"`
}
