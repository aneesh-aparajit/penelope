package config

type AgentDef struct {
	CLI        CLI      `json:"cli"`
	Model      Model    `json:"model"`
	Role       string   `json:"role"`
	Soul       string   `json:"soul"`
	CanMessage []string `json:"canMessage"`
	MaxRetries uint8    `json:"maxRetries"`
}
