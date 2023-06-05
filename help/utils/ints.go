package utils

import "strconv"

// IntToString
//  @Description:   int 转 string
//  @param e
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:26:40
func IntToString(e int) string {

	return strconv.Itoa(e)
}

// Int64ToString
//  @Description: int64 转 string
//  @param e
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:27:37
func Int64ToString(e int64) string {
	return strconv.FormatInt(e, 10)
}
