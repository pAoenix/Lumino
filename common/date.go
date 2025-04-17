package common

import "time"

// GetFirstDayOfMonth 获取指定时间的月份第一天
func GetFirstDayOfMonth(t time.Time) (firstDay time.Time) {
	year, month, _ := t.Date()

	// 本月第一天
	firstDay = time.Date(year, month, 1, 0, 0, 0, 0, t.Location())

	return firstDay
}

// GetLastDayOfMonth 获取指定时间的月份最后一天
func GetLastDayOfMonth(t time.Time) (lastDay time.Time) {
	year, month, _ := t.Date()

	// 本月第一天
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())

	// 本月最后一天
	lastDay = firstDay.AddDate(0, 1, -1)

	return lastDay
}

// 使用示例
/*
now := time.Now()
first, last := GetFirstAndLastDayOfMonth(now)
fmt.Println("First day:", first.Format("2006-01-02"))
fmt.Println("Last day:", last.Format("2006-01-02"))
*/
