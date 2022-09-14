# 打印素数

## 问题描述

每行5个输出从1到20000之间的所有素数。

## 解决方案

### 循环

```go
func main() {
    var count int
    for i := 2; i <= 20000; i++ {
        if isPrime(i) {
            fmt.Printf("%d\t", i)
            count++
            if count%5 == 0 {
                fmt.Println()
            }
        }
    }
}

func isPrime(n int) bool {
    for i := 2; i < n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}
```


## 优化分析

最费时的是判断素数，可以使用开方来减少循环次数。


### 优化

```go
func main() {
    var count int
    for i := 2; i <= 20000; i++ {
        if isPrime(i) {
            fmt.Printf("%d\t", i)
            count++
            if count%5 == 0 {
                fmt.Println()
            }
        }
    }
}

func isPrime(n int) bool {
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}
```

#### 性能测试
    
    ```go 
    func BenchmarkIsPrime(b *testing.B) {
        for i := 0; i < b.N; i++ {
            isPrime(20000)
        }
    }
    ```
##### 优化前
    
```shell
% go test -bench=.  
goos: darwin
goarch: arm64
pkg: soft-project/task1/03print_nums
BenchmarkIsPrime-10     1000000000               0.6320 ns/op
PASS
ok      soft-project/task1/03print_nums 1.096s
```

##### 优化后

```shell
% go test -bench=.
goos: darwin
goarch: arm64
pkg: soft-project/task1/03print_nums
BenchmarkIsPrime-10     1000000000               0.6288 ns/op
PASS
ok      soft-project/task1/03print_nums 0.775s
```


## 总结

优化后测试时间减少了0.3s，大约百分之三十。可见优化的效果还是很明显的。



