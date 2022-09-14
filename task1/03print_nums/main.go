package main

import "fmt"

// 每行5个打印1-20000的素数
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

// 判断是否为素数
func isPrime(n int) bool {
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
