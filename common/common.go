package common

// ContainsInt 判断一个int是否在int数组中
func ContainsInt(arr []int, target int) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}
