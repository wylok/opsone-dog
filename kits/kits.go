package kits

import (
	"agent-dog/config"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/golang-module/carbon"
	"github.com/mitchellh/go-ps"
	"os"
	"strconv"
)

func WritePid() {
	_ = fileutil.WriteStringToFile(config.AgentDogPidFile, strconv.Itoa(os.Getpid()), false)
}
func WriteLog(msg string) {
	logfile := config.LogFile + "." + carbon.Now().ToDateString()
	if fileutil.IsExist(logfile) {
		_ = fileutil.WriteStringToFile(logfile, carbon.Now().ToDateTimeString()+" opsone-dog "+fmt.Sprintln(msg), true)
	} else {
		file, _ := os.Create(logfile)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
}
func CheckAgentPid() bool {
	f, _ := fileutil.ReadFileToString(config.AgentPidFile)
	if f != "" {
		v, err := strconv.Atoi(f)
		if err == nil {
			p, err := ps.FindProcess(v)
			if err == nil && p != nil {
				return true
			}
		}
	}
	return false
}
func CheckAgentDogPid() {
	f, _ := fileutil.ReadFileToString(config.AgentDogPidFile)
	if f != "" {
		v, err := strconv.Atoi(f)
		if err == nil {
			p, err := ps.FindProcess(v)
			if err == nil && p != nil {
				if p.Pid() != os.Getpid() {
					config.Qch <- "进程(" + strconv.Itoa(os.Getpid()) + ")其它dog进程正在运行,退出"
				}
			} else {
				WritePid()
			}
		}
	} else {
		WritePid()
	}
}
