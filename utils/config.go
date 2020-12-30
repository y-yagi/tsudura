package utils

type Config struct {
	Root     string `toml:"root"`
	Endpoint string `toml:"endpoint"`
	Bucket   string `toml:"bucket"`
	Region   string `toml:"region"`
	Secret   string `toml:"secret"`
	Token    string `toml:"token"`
}