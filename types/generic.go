package types

type SQLCredentials struct {
	Host         string `toml:"host"`
	User         string `toml:"user"`
	Port         int    `toml:"port"`
	Pass         string `toml:"pass"`
	DatabaseName string `toml:"db-name"`
}
