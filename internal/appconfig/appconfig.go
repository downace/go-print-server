package appconfig

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	CertFile string `yaml:"certFile" json:"certFile"`
	KeyFile  string `yaml:"keyFile" json:"keyFile"`
}

type AppConfig struct {
	Host            string            `yaml:"host" json:"host"`
	Port            uint16            `yaml:"port" json:"port"`
	ResponseHeaders map[string]string `yaml:"responseHeaders" json:"responseHeaders"`
	TLS             TLSConfig         `yaml:"tls" json:"tls"`
}
