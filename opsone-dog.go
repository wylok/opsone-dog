package main

import (
	"agent-dog/config"
	_ "agent-dog/daemon"
	"agent-dog/kits"
	"agent-dog/task"
	"bufio"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/netutil"
	"net"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	go func() {
		_ = os.MkdirAll(config.LogPath, 0755)
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						kits.WriteLog(fmt.Sprint(r))
					}
				}()
				kits.CheckAgentDogPid()
				f := config.AgentPath + "/config.ini"
				if fileutil.IsExist(f) {
					remoteIp, _ := fileutil.ReadFileToString(f)
					remoteIp = strings.TrimSpace(remoteIp)
					Ip := remoteIp
					if strings.Contains(remoteIp, ":") {
						Ip = strings.Split(remoteIp, ":")[0]
					}
					if netutil.IsPingConnected(Ip) {
						config.ServerAddr = "http://" + remoteIp
					} else {
						kits.WriteLog("无效服务端连接地址:" + remoteIp)
					}
				}
			}()
			time.Sleep(30 * time.Second)
		}
	}()
	go func() {
		listen, err := net.Listen("tcp", "127.0.0.1:54321")
		if err == nil {
			for {
				conn, _ := listen.Accept()
				var (
					buf [128]byte
					ok  bool
				)
				reader := bufio.NewReader(conn)
				n, err := reader.Read(buf[:])
				if err == nil {
					_, _ = conn.Write([]byte(config.Version))
					msg := string(buf[:n])
					if msg == "uninstall" {
						ok = true
					}
					if strings.Contains(msg, "upgrade:") {
						if strings.Split(msg, ":")[1] != config.Version {
							ok = true
						}
					}
					if ok {
						for _, f := range []string{config.AgentDogFile, config.AgentDogPidFile} {
							_ = fileutil.RemoveFile(f)
						}
						config.Qch <- "进程(" + strconv.Itoa(os.Getpid()) + ")版本更新,已退出"
					}
				}
			}
		}
	}()
}
func main() {
	go kits.WritePid()
	go task.Scheduler()
	select {
	case msg := <-config.Qch:
		kits.WriteLog(msg)
		os.Exit(0) //进程异常结束
	}
}
