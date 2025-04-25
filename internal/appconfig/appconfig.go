package appconfig

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	CertFile string `yaml:"certFile" json:"certFile"`
	KeyFile  string `yaml:"keyFile" json:"keyFile"`
}

type AuthConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type AppConfig struct {
	Host            string            `yaml:"host" json:"host"`
	Port            uint16            `yaml:"port" json:"port"`
	ResponseHeaders map[string]string `yaml:"responseHeaders" json:"responseHeaders"`
	TLS             TLSConfig         `yaml:"tls" json:"tls"`
	Auth            AuthConfig        `yaml:"auth" json:"auth"`
}

func NewDefaultConfig() AppConfig {
	return AppConfig{
		Host:            "0.0.0.0",
		Port:            8888,
		ResponseHeaders: map[string]string{},
	}
}
