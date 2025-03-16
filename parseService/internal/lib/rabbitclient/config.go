package rabbitclient

type Config struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Channel  string `yaml:"channel"`
}
