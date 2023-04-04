package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// StringToInt
//  @Description: string 转 int
//  @param s
//  @return int
//  @Author  ahKevinXy
//  @Date2023-04-04 14:29:33
func StringToInt(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}

// Substr2
//  @Description:  截取文件
//  @param str
//  @param start
//  @param end
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:48:04
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

// StringToInt64
//  @Description:   String 转 Int64
//  @param s
//  @return int64
//  @Author  ahKevinXy
//  @Date2023-04-04 14:30:21
func StringToInt64(s string) int64 {
	r, _ := strconv.ParseInt(s, 10, 64)
	return r
}

// UpperFirst
//  @Description:   首字母大写
//  @param s
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:32:18
func UpperFirst(s string) string {
	if len(s) > 0 {
		strings.Replace(s, s[0:1], strings.ToUpper(s[0:1]), 1)
	}

	return s
}

// HideStar
// @Description: 字符串处理函数
// @Author ahKevinXy
// @Date 2022-11-09 19:39:58
func HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str, "@") {
		// 邮箱
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := Substr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			// 手机号
			result = Substr2(str, 0, 3) + "****" + Substr2(str, 7, 11)
		} else {
			nameRune := []rune(str)
			lens := len(nameRune)
			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 {
				result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
			}
		}
		return
	}
}

// MaskMobile
// @Description: 隐藏手机号码
// @Author ahKevinXy
// @Date 2022-11-09 19:40:59
func MaskMobile(mobile string, maskNumber int) string {
	offset := 3
	if len(mobile) > 11 {
		offset += 2
	}
	maskNumber += offset

	ret := mobile
	if len(mobile) > maskNumber {
		_cardNo := make([]rune, len(mobile))
		for index, _c := range mobile {
			if index >= offset && index < maskNumber {
				_cardNo[index] = '*'
				continue
			}
			_cardNo[index] = _c
		}

		ret = string(_cardNo)
	}

	return ret
}

// MaskCustomerName
// @Description: 隐藏用户名
// @Author ahKevinXy
// @Date 2022-11-09 19:41:09
func MaskCustomerName(name string) string {

	_name := []rune(name)
	if len(name) > 4 {
		name = string(_name[:2]) + "***" + string(_name[len(_name)-2:])
	} else if len(name) < 1 {
		name = ""
	} else if len(name) <= 4 {
		name = string(_name[:1]) + "***"
	}

	return name
}

// MaskIDCard
// @Description:  隐藏身份证号码
// @Author ahKevinXy
// @Date 2022-11-09 19:41:18
func MaskIDCard(idcard string, maskNumber int) string {
	offset := 4
	maskNumber += offset

	ret := idcard
	if len(idcard) > maskNumber {
		_cardNo := make([]rune, len(idcard))
		for index, _c := range idcard {
			if index >= offset && index < maskNumber {
				_cardNo[index] = '*'
				continue
			}
			_cardNo[index] = _c
		}

		ret = string(_cardNo)
	}

	return ret
}

// MaskBackCard
// @Description:  隐藏银行卡
// @Author ahKevinXy
// @Date 2022-11-09 19:41:25
func MaskBackCard(bankcard string, maskNumber int) string {
	offset := 4
	maskNumber += offset
	if len(bankcard) > 16 {
		maskNumber += len(bankcard) - 16
	}

	ret := bankcard
	if len(bankcard) > maskNumber {
		_cardNo := make([]rune, len(bankcard))
		for index, _c := range bankcard {
			if index >= offset && index < maskNumber {
				_cardNo[index] = '*'
				continue
			}
			_cardNo[index] = _c
		}

		ret = string(_cardNo)
	}

	return ret
}
