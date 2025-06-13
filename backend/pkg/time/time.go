package time

import (
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
	TimeZone   = "Asia/Shanghai"
)

// 格式化时间（支持时区，默认使用Asia/Shanghai）
func FormatTime(t time.Time, tz ...string) string {
	timeZone := TimeZone
	if len(tz) > 0 {
		timeZone = tz[0]
	}
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return t.Format(TimeLayout)
	}
	return t.In(location).Format(TimeLayout)
}

// 解析时间字符串（支持时区，默认使用Asia/Shanghai）
func ParseTime(timeStr string, tz ...string) (time.Time, error) {
	timeZone := TimeZone
	if len(tz) > 0 {
		timeZone = tz[0]
	}
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Parse(TimeLayout, timeStr)
	}
	return time.ParseInLocation(TimeLayout, timeStr, location)
}
