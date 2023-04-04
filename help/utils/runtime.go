package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentDir
//  @Description:   获取当前运行路径
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:41:44
func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
