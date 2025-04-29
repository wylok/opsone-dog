package config

const (
	Version         = "2025033102"
	AgentPath       = "/opt/opsone"
	PidPath         = "/var/run"
	LogPath         = "/var/log/opsone"
	AgentFile       = AgentPath + "/opsone-agent"
	AgentDogFile    = AgentPath + "/opsone-dog"
	AgentPidFile    = PidPath + "/opsone-agent.pid"
	AgentDogPidFile = PidPath + "/opsone-dog.pid"
	LogFile         = LogPath + "/opsone-agent.log"
)

var (
	Qch        = make(chan string, 1)
	ServerAddr string
)
