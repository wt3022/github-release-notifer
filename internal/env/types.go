package env

type EmailConfig struct {
	SMTPHost  string
	SMTPPort  string
	FromEmail string
}

type Config struct {
	Env            string
	GithubToken    string
	ProductionMode bool
	EmailConfig    EmailConfig
}
