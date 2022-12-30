package finetuneimpl

type configSetDefault interface {
	setDefault()
}

type configValidate interface {
	validate() error
}

type Config struct {
	Modelarts ModelartsConfig `json:"modelarts"   required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Modelarts,
	}
}

func (cfg *Config) Validate() error {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configValidate); ok {
			if err := v.validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configSetDefault); ok {
			v.setDefault()
		}
	}
}

type ModelartsConfig struct {
	AccessKey   string `json:"access_key" required:"true"`
	SecretKey   string `json:"secret_key" required:"true"`
	Region      string `json:"region" required:"true"`
	ProjectName string `json:"project_name" required:"true"`
	ProjectId   string `json:"project_id" required:"true"`

	// modelarts endpoint
	Endpoint string `json:"endpoint" required:"true"`
}
