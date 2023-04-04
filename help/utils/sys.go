package utils

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// Reload
//  @Description:   重载
//  @return error
//  @Author  ahKevinXy
//  @Date2023-04-04 14:41:19
func Reload() error {
	if runtime.GOOS == "windows" {
		return errors.New("system not allow")
	}
	pid := os.Getppid()

	cmd := exec.Command("kill", "-1", strconv.Itoa(pid))

	return cmd.Run()
}
