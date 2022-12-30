package domain

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxFinetuneNameLength int `json:"max_finetune_name_length"`
	MinFinetuneNameLength int `json:"min_finetune_name_length"`

	// Key is the finetue model name
	Finetunes map[string]FinetuneParameterConfig `json:"finetunes" required:"true"`
}

func (r *Config) SetDefault() {
	if r.MaxFinetuneNameLength == 0 {
		r.MaxFinetuneNameLength = 50
	}

	if r.MinFinetuneNameLength == 0 {
		r.MinFinetuneNameLength = 5
	}
}

type FinetuneParameterConfig struct {
	Tasks           []string `json:"tasks"           required:"true"`
	Hyperparameters []string `json:"hyperparameters" required:"true"`
}
