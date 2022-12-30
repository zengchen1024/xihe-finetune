package config

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-finetune/domain"
	"github.com/opensourceways/xihe-finetune/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-finetune/infrastructure/watchimpl"
)

type configSetDefault interface {
	SetDefault()
}

type configValidate interface {
	Validate() error
}

type Config struct {
	Watch    watchimpl.Config    `json:"watch"        required:"true"`
	Domain   domain.Config       `json:"domain"       required:"true"`
	Finetune finetuneimpl.Config `json:"finetune"     required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Watch,
		&cfg.Domain,
		&cfg.Finetune,
	}
}

func (cfg *Config) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configValidate); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *Config) setDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configSetDefault); ok {
			v.SetDefault()
		}
	}
}

func (cfg *Config) InitDomain() {
	domain.Init(&cfg.Domain)
}

func LoadConfig(path string) (*Config, error) {
	v := new(Config)

	if err := utils.LoadFromYaml(path, v); err != nil {
		return nil, err
	}

	v.setDefault()

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}
