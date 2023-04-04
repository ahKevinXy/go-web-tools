package utils

import "time"

const (
	GoLangTimeFormat      = "2006-01-02 15:04:05"
	GoLangPointTimeFormat = "2006.01.02"
	GolandDayTimeFormat   = "2006-01-02"
	GolandDayTimePath     = "2006/01/02"
)

const (
	TimeOffset       = 8 * 3600  //8 hour offset
	TimeHalfOffset   = 12 * 3600 //Half-day hourly offset
	TimeDayOffset    = 24 * 3600
	TimeHourOffset   = 3600
	TimeMinuteOffest = 60
)

// GetTimeDir
//  @Description:   获取时间路径
//  @return string
//  @Author  ahKevinXy
//  @Date2023-04-04 14:44:28
func GetTimeDir() string {
	n := time.Now()

	return n.Format(GolandDayTimePath)
}
