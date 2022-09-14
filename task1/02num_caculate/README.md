# 找出数组中子数组的最大和

## 问题描述

给定一个数组，找出数组中子数组的最大和。

## 解决方案

### 动态规划

动态规划的思想是，如果前面的子数组的和是负数，那么就不要它，从当前位置开始重新计算子数组的和。

```go
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
```

## 正确性证明

```go
func TestMaxSubArray(t *testing.T) {
var arr = []int{1, -2, 3, 10, -4, 7, 2, -5}
var max = 18
var result = maxSubArray(arr)
if result != max {
t.Errorf("maxSubArray(%v) = %d; want %d", arr, result, max)
    }
}
```

```shell
$ go test -v
=== RUN   TestMaxSubArray
--- PASS: TestMaxSubArray (0.00s)
PASS
ok      soft-project/task1/02num_caculate       0.262s
```

## 性能测试


```go
func BenchmarkMaxSubArray(b *testing.B) {
    var arr = []int{1, -2, 3, 10, -4, 7, 2, -5}
    for i := 0; i < b.N; i++ {
        maxSubArray(arr)
    }
}
```

```shell
$ go test -bench=.  
goos: darwin
goarch: arm64
pkg: soft-project/task1/02num_caculate
BenchmarkMaxSubArray-10         230510385                5.066 ns/op
PASS
ok      soft-project/task1/02num_caculate       2.258s
```
