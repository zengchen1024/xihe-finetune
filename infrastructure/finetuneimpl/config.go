package finetuneimpl

type Config struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`

	// finetune endpoint
	Endpoint string `json:"endpoint" required:"true"`
}
