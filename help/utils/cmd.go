package utils

import (
	"bytes"
	"errors"
	"os/exec"
)

// ShellExec
// @Description: 在Shell执行命令
// @Author ahKevinXy
// @Date 2022-12-01 11:44:59
func ShellExec(cmd string) (*bytes.Buffer, error) {
	command := exec.Command("sh")
	in := bytes.NewBuffer(nil)
	out := bytes.NewBuffer(nil)
	errOut := bytes.NewBuffer(nil)
	command.Stdin = in
	command.Stdout = out
	command.Stderr = errOut
	in.WriteString(cmd)
	in.WriteString("\n")
	in.WriteString("exit\n")
	if err := command.Run(); err != nil {
		return nil, errors.New(errOut.String())
	}
	return out, nil
}
