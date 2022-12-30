package watchimpl

type Config struct {
	// Interval specifies the interval of second between two loops
	// that check all trainings in a loop.
	Interval int `json:"interval"`

	Endpoint    string `json:"endpoint"      required:"true"`
	MaxWatchNum int    `json:"max_watch_num" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}
}
