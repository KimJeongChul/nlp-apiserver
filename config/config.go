package config

// ServerConfigJson Server config
type ServerConfigJson struct {
	ListenPort  string `json:"listenPort"`
	HomePath    string `json:"homePath"`
	LogPath     string `json:"logPath"`
	LogLevel    string `json:"LogLevel"`
	Ssl         int    `json:"ssl"`
	CertPemPath string `json:"certPemPath"`
	KeyPemPath  string `json:"keyPemPath"`
}
