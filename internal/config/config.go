package config

type (
	RdmConfig struct {
		Commands map[string]userCommand `json:"commands"`
	}

	userCommand struct {
		ExecutablePath string `json:"executablePath"`
	}
)

func New() *RdmConfig {
	return &RdmConfig{Commands: map[string]userCommand{}}
}
