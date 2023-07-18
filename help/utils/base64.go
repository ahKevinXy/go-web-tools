package utils

import "encoding/base64"

// Base64Encode
// @Description: base64加密
// @Author ahKevinXy
// @Date 2022-11-08 16:49:08
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64DecodeToString
// @Description: base64解码
// @Author ahKevinXy
// @Date 2022-11-08 16:49:19
func Base64DecodeToString(data string) string {
	decodedByte, _ := base64.StdEncoding.DecodeString(data)
	return string(decodedByte)
}

// Base64Decode
// @Description: base64解码
// @Author ahKevinXy
// @Date 2022-11-22 17:04:52
func Base64Decode(data string) []byte {
	decodedByte, _ := base64.StdEncoding.DecodeString(data)
	return decodedByte
}

func Base64DecodeToByte(data string) []byte {
	decodedByte, _ := base64.StdEncoding.DecodeString(data)
	return decodedByte
}
// Base64EncodeToByte
// 添加 byte 数据
func Base64EncodeToByte(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
