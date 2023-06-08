package cache

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     "6379",
		Username: "",
		Password: "",
		DB:       0,
	}
}
