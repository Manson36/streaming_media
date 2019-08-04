package main

import "fmt"

func quickSort(values []int, left, right int) {
	temp := values[left]
	i, j := left, right
	p := left

	for i <= j {
		for j >= p && values[j] >= temp {
			j--
		}
		if j >= p {
			values[p] = values[j]
			p = j
		}

		for i <= p && values[i] <= temp {
			i++
		}
		if i <= p {
			values[p] = values[i]
			p = i
		}
	}

	values[p] = temp
	if p - left > 1 {
		quickSort(values, left, p)
	}

	if right - p > 1 {
		quickSort(values, p + 1, right)
	}
}

func QuickSort(values []int) {
	if len(values) <= 1 {
		return
	}

	quickSort(values, 0, len(values) - 1 )
}

func QuickSort2(values []int) {
	
}

func main() {
	var values = []int{22,3,4,55,6,24,11}
	QuickSort(values)
	fmt.Println(values)
}
