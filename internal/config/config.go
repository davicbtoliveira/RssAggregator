package config

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg Config) SetUser() error {
	cfg.CurrentUserName = "dcbto"

	if err := write(cfg); err != nil {
		return err
	}

	return nil
}
