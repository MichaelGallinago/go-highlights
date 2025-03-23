package grpcserver

type Config struct {
	Network string `yaml:"network"`
	Port    int    `yaml:"port"`
}
