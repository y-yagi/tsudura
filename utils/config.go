package utils

type Config struct {
	Root     string `toml:"root"`
	Endpoint string `toml:"endpoint"`
	Bucket   string `toml:"bucket"`
	Region   string `toml:"region"`
	Secret   string `toml:"secret"`
	Token    string `toml:"token"`
	AddOnly  bool   `toml:"addonly"`
	Debug    bool   `toml:"debug"`
}
