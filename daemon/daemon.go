package daemon

import (
	"agent-dog/config"
	"github.com/duke-git/lancet/v2/fileutil"
	"os"
	"os/exec"
	"strconv"
)

func init() {
	RunCheck()
	if len(os.Args) >= 2 {
		cmd := exec.Command(os.Args[0])
		_ = cmd.Start()
		os.Exit(0)
	}
}
func RunCheck() {
	// 检查是否重复启动
	if fileutil.IsExist(config.AgentDogPidFile) {
		f, err := fileutil.ReadFileToString(config.AgentDogPidFile)
		if err == nil && f != "" {
			v, _ := strconv.Atoi(f)
			if v != os.Getpid() {
				_ = fileutil.RemoveFile(config.AgentDogPidFile)
				config.Qch <- "进程(" + strconv.Itoa(os.Getpid()) + ")重复启动,退出"
			}
		}
	}
}
