package chessboard

import "fmt"

func PrintChessboard() {
	n, m := readChesboardSize()

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if (i+j)%2 == 0 {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func readChesboardSize() (int, int) {
	var n int
	var m int
	fmt.Printf("Enter the width of the board: ")
	fmt.Scan(&n)
	fmt.Printf("Enter the length of the board: ")
	fmt.Scan(&m)

	return n, m
}
