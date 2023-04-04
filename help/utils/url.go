package utils

import (
	"fmt"
	"strings"
)

// BuildUrl
//  @Description:  创建请求链接
//  @param prefix
//  @param params
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:45:12
func BuildUrl(prefix string, params ...interface{}) string {
	var (
		l       int
		linkUrl []string
	)

	linkUrl = append(linkUrl, "/"+strings.Trim(prefix, "/"))
	l = len(params)
	if l != (l/2)*2 {
		l = (l / 2) * 2
	}
	if l > 0 {
		for i := 0; i < l; {
			k := fmt.Sprintf("%v", params[i])
			v := fmt.Sprintf("%v", params[i+1])
			if len(k) > 0 && v != "0" {
				linkUrl = append(linkUrl, fmt.Sprintf("%v/%v", k, v))
			}
			i += 2
		}
	}
	return strings.TrimRight(strings.Join(linkUrl, "/"), "/")
}
