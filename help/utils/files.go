package utils

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// PathExists
//  @Description:  判断文件是否存在
//  @param path
//  @return bool
//  @return error
//  @Author  ahKevinXy
//  @Date2023-04-04 14:53:21
func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("have same file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDir
// @Description: 创建目录
// @Date 2022-11-08 16:43:27
func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {

			if err := os.MkdirAll(v, os.ModePerm); err != nil {

				return err
			}
		}
	}
	return err
}

// FileMove
// @Description: 移动目录
// @Auth ahKevinXy
// @Date 2022-11-08 16:45:04
func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	revoke := false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

// DeLFile
// @Description: 删除文件
// @Author ahKevinXy
// @Date 2022-11-08 16:46:11
func DeLFile(filePath string) error {
	return os.RemoveAll(filePath)
}

// FileExist
// @Description: 判断文件是否存在
// @Author ahKevinXy
// @Date 2022-11-08 16:46:24
func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

// FileMd5
// @Description: 获取文件 md5
// @Author ahKevinXy
// @Date 2022-11-08 16:46:34
func FileMd5(path string) (string, error) {
	var md5str string
	var err error
	var file *os.File
	file, err = os.Open(path)

	if err != nil {
		return md5str, err
	}
	defer file.Close()
	md5h := md5.New()
	_, err = io.Copy(md5h, file)
	if err == nil {
		md5str = fmt.Sprintf("%x", md5h.Sum(nil))
	}
	return md5str, err
}

// IsImage
// @Description: 判断是否是图片
// @Author ahKevinXy
// @Date 2022-11-08 16:46:41
func IsImage(path string) bool {
	slice := strings.Split(path, ".")
	ext := strings.ToLower(strings.TrimSpace(slice[len(slice)-1]))
	exts := map[string]string{"jpeg": "jpeg", "jpg": "jpg", "gif": "gif", "png": "png", "bmp": "bmp", "tif": "tif", "tiff": "tiff"}
	_, ok := exts[ext]
	return ok
}

// GetSuffixName
// @Description: 获取文件的后缀名
// @Author ahKevinXy
// @Date 2022-12-01 11:49:09
func GetSuffixName(str, seg string) string {
	slice := strings.Split(str, seg)
	l := len(slice)
	if l > 1 {
		return slice[(l - 1)]
	}
	return ""
}

//FormatByte 转换字节大小
func FormatByte(size int) string {
	fSize := float64(32)
	units := [6]string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i := 0; fSize > 1024 && i < 5; i++ {
		fSize /= 1024
	}
	num := fmt.Sprintf("%.2f", fSize)
	return string(num) + " " + units[i]
}

// ScanDir
// @Description: ScanDir扫描目录目录中的文件
// @Author ahKevinXy
// @Date 2022-12-01 11:48:20
func ScanDir(dir string) (files []string, err error) {
	dir = strings.TrimSuffix(dir, "/")
	if infos, err := ioutil.ReadDir(dir); err == nil {
		for _, info := range infos {
			file := dir + "/" + info.Name()
			if info.IsDir() {
				item, err := ScanDir(file)
				if err != nil {
					return nil, err
				}
				files = append(files, item...)
			} else {
				files = append(files, file)
			}
		}
	} else {
		return nil, err
	}
	return
}

// CopyFile
// @Description: 复制文件
// @Author ahKevinXy
// @Date 2022-12-01 11:48:59
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer dst.Close()
	return io.Copy(dst, src)
}

// ImageFileToBase64  图片文件转 base64
func ImageFileToBase64(localFilePath string) string {

	// 读取文件
	f, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(f)
}

// IsDir
// @Description: 判断是不是路径
// @Author ahKevinXy
// @Date 2022-12-01 11:49:17
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile
// @Description: 判断是否是文件
// @Author ahKevinXy
// @Date 2022-12-01 11:47:47
func IsFile(path string) bool {
	return !IsDir(path)
}

// MkDir
// @Description: 安装递归 创建目录
// @Author ahKevinXy
// @Date 2022-12-01 11:47:56
func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// IsTextFile returns true if file content format is plain text or empty.
func IsTextFile(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return strings.Contains(http.DetectContentType(data), "text/")
}

// IsImageFile
// @Description: 是否是文件数据
// @Author ahKevinXy
// @Date 2023-01-10 18:02:50
func IsImageFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "image/")
}

// IsPDFFile
// @Description: 是否是PDF文件
// @Author ahKevinXy
// @Date 2023-01-10 18:02:59
func IsPDFFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "application/pdf")
}

// IsVideoFile
// @Description: 是否是视频文件
// @Author ahKevinXy
// @Date 2023-01-10 18:03:10
func IsVideoFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "video/")
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

// IsSameSiteURLPath returns true if the URL path belongs to the same site, false otherwise.
// False: //url, http://url, /\url
// True: /url
func IsSameSiteURLPath(url string) bool {
	return len(url) >= 2 && url[0] == '/' && url[1] != '/' && url[1] != '\\'
}

func IsMaliciousPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "..")
}

// FileIsExisted returns     bool
// @Description: 判断文件是否存在
// @Author ahKevinXy
// @Date 2023-01-11 11:46:41
func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		existed = false
	}
	return existed
}
