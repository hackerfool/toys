package main

import "fmt"

func main() {
	array := []int{-1, 1, 2, 3, 4, 5, 6, 7}
	for x := 0; x <= 10; x++ {
		pos, err := binarySearch(array, x, 0, len(array)-1)
		if err != true {
			fmt.Println("not found : ", x)
		} else {
			fmt.Println("pos : ", pos, " value : ", x)
		}
	}
}

func binarySearch(array []int, value int, start int, end int) (int, bool) {
	for start <= end {
		mid := start + int((end-start)/2)
		// fmt.Println(start, ",", mid, ",", end)
		if array[mid] == value {
			return mid, true
		}
		if array[mid] > value {
			end = mid - 1
		}
		if array[mid] < value {
			start = mid + 1
		}
	}
	return -1, false
}
