package main

import "fmt"

func getPrime(number int, primes []int) []int {
	if number == 1 {
		// fmt.Println(&primes)
		primes = append(primes, number)
		// fmt.Println(&primes)
		return primes
	}
	for _, value := range primes {
		if number%value == 0 && value != 1 {
			return primes
		}
	}
	primes = append(primes, number)
	return primes
}

func main() {
	primes := make([]int, 0, 100)
	for i := 1; i < 100; i++ {
		primes = getPrime(i, primes)
	}
	for _, value := range primes {
		fmt.Println(value)
	}
}
