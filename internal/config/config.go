package config

type Config struct {
	AppConfig appConfig `json:"app"`
}

type appConfig struct {
	Cache bool `json:"cache"`
}
