package db

type Config struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"db_name" yaml:"db_name"`
}

func DefaultConfig() Config {
	return Config{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "admin",
		Password: "123456",
		DBName:   "zura",
	}
}
