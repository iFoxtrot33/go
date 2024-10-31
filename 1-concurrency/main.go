package main

import (
	"fmt"
	"math/rand"
)

func main() {
	numbers := make(chan int)
	squares := make(chan int)

	go generateRandomSlice(numbers)
	go makeSquare(numbers, squares)

	results := make([]int, 0, 10)

	for i := 0; i < 10; i++ {
		num := <-squares
		results = append(results, num)
	}
	fmt.Print("Результаты: ", results)
}

func generateRandomSlice(out chan int) {
	for i := 0; i < 10; i++ {
		num := rand.Intn(101)
		out <- num

	}
	close(out)
}

func makeSquare(in chan int, out chan int) {
	for num := range in {
		square := num * num
		out <- square
	}
	close(out)
}
