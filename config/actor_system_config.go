package config


type ActorSystemConfig struct {
	Minactor int `env:"minactor" default:"10"`
	Maxactor int `env:"maxactor" default:"100"`
	AutoScale
}

type AutoScale struct {
	UpscaleQueueSize int `env:"upscalequeuesize" default:"100"`
	DownscaleQueueSize int `env:"downscalequeuesize" default:"10"`
}