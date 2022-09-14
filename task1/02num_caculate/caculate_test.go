package main

import "testing"

// 正确性测试
func TestMaxSubArray(t *testing.T) {
	var arr = []int{1, -2, 3, 10, -4, 7, 2, -5}
	var max = 18
	var result = maxSubArray(arr)
	if result != max {
		t.Errorf("maxSubArray(%v) = %d; want %d", arr, result, max)
	}
}

// 性能测试
func BenchmarkMaxSubArray(b *testing.B) {
	var arr = []int{1, -2, 3, 10, -4, 7, 2, -5}
	for i := 0; i < b.N; i++ {
		maxSubArray(arr)
	}
}
