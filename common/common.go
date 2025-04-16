package common

import (
	"fmt"
	"github.com/lib/pq"
	"strconv"
)

// ContainsInt 判断一个int是否在int数组中
func ContainsInt(arr []int, target int) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

// ContainsUint 判断一个int是否在int数组中
func ContainsUint(arr []uint, target uint) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

// String2Uint -
func String2Uint(str string) (uint, error) {
	u64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return 0, err
	}

	// 2. 将uint64转换为uint
	u := uint(u64)
	return u, nil
}

// ConvertArrayToIntSlice -
func ConvertArrayToIntSlice(int32Array pq.Int32Array) []int {
	intSlice := make([]int, len(int32Array))
	for i, v := range int32Array {
		intSlice[i] = int(v)
	}
	return intSlice
}
