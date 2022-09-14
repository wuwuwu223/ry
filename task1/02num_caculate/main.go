package main

import "fmt"

// 找出数组中子数组的最大和
func main() {
	var n int
	fmt.Printf("请输入数组长度:")
	fmt.Scanf("%d", &n)
	var arr = make([]int, n)
	for j := 0; j < n; j++ {
		fmt.Printf("请输入第%d个数:", j+1)
		fmt.Scanf("%d", &arr[j])
	}
	result := maxSubArray(arr)
	fmt.Println(result)

}

func maxSubArray(nums []int) int {
	max := nums[0]
	sum := 0
	for i := 0; i < len(nums); i++ {
		if sum > 0 {
			sum += nums[i]
		} else {
			sum = nums[i]
		}
		if sum > max {
			max = sum
		}
	}
	return max
}
