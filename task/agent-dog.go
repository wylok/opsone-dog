package task

import (
	"agent-dog/config"
	"agent-dog/kits"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/duke-git/lancet/v2/system"
	"github.com/jakecoffman/cron"
	"os"
)

func Scheduler() {
	c := cron.New()
	c.AddFunc("0 * * * * *", CheckAgent, "CheckAgent")
	c.Start()
}
func CheckAgent() {
	if !fileutil.IsExist(config.AgentPath + "/config.ini") {
		_ = netutil.DownloadFile(config.AgentPath+"/config.ini", config.ServerAddr+"/api/v1/conf/config.ini")
	}
	if !fileutil.IsExist(config.AgentFile) {
		kits.WriteLog("下载agnet:" + config.ServerAddr + "/api/v1/ag/opsone-agent")
		_ = netutil.DownloadFile(config.AgentFile, config.ServerAddr+"/api/v1/ag/opsone-agent")
	}
	if fileutil.IsExist(config.AgentFile) {
		_ = os.Chmod(config.AgentFile, 0755)
		if fileutil.IsExist(config.AgentPidFile) {
			if !kits.CheckAgentPid() {
				_, stderr, err := system.ExecCommand(config.AgentFile + " start")
				if err != nil || stderr != "" {
					_ = fileutil.RemoveFile(config.AgentFile)
				}
			}
		} else {
			_, stderr, err := system.ExecCommand(config.AgentFile + " start")
			if err != nil || stderr != "" {
				_ = fileutil.RemoveFile(config.AgentFile)
			}
		}
	}
}
