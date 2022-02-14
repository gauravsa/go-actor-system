package actor


type Config struct {
	MinActor int `env:"min_actor" default:"10"`
	MaxActor int `env:"max_actor" default:"100"`
	AutoScale
}

type AutoScale struct {
	UpscaleQueueSize int `env:"upscale_queue_size" default:"100"`
	DownscaleQueueSize int `env:"downscale_queue_size" default:"10"`
}