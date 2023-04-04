package utils

import (
	"math/rand"
	"time"
)

// RandStr
//  @Description:  获取随机数
//  @param size
//  @param kind
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:37:39
func RandStr(size int, kind int) string {

	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)

	isAll := kind > 2 || kind < 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]

		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
