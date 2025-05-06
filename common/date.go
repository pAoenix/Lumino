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

// CheckDailyCoverage 检查日期列表是否包含时间范围内的所有天数
func CheckDailyCoverage(beginDate, endDate time.Time, dateList []string) (bool, []string) {
	var missingDates []string

	// 创建日期集合便于快速查找
	dateSet := make(map[string]struct{})
	for _, date := range dateList {
		dateSet[date] = struct{}{}
	}

	// 遍历时间范围内的每一天
	for d := beginDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if _, exists := dateSet[dateStr]; !exists {
			missingDates = append(missingDates, dateStr)
		}
	}

	return len(missingDates) == 0, missingDates
}

// CheckMonthlyCoverage 检查日期列表是否包含时间范围内的所有月份
func CheckMonthlyCoverage(beginDate, endDate time.Time, dateList []string) (bool, []string) {
	var missingMonths []string

	// 创建月份集合便于快速查找
	monthSet := make(map[string]struct{})
	for _, date := range dateList {
		monthSet[date] = struct{}{}
	}

	// 遍历时间范围内的每个月
	for d := beginDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
		monthStr := d.Format("2006-01")
		if _, exists := monthSet[monthStr]; !exists {
			missingMonths = append(missingMonths, monthStr)
		}
	}

	return len(missingMonths) == 0, missingMonths
}
