package main

import "fmt"

// Global varaibles to store position of last inserted token
var global_row, global_column int

// Prints the board to terminal
func print_board(board [][]int) {

	fmt.Println("")

	// Cycles threw the full board, printing
	// X for blank, and colored 0's for players
	for n := 0; n < len(board); n++ {

		for m := 0; m < len(board[0]); m++ {

			// Switch to compare value. Faster than if statements
			// and no breaks needed in google go switch
			// Colors depend on support from terminal!
			switch board[n][m] {
			case -1:
				fmt.Print("X ")

			case 0:
				fmt.Printf("\x1B")
				fmt.Printf("[36m")
				fmt.Print("0 ")
				fmt.Printf("\033")
				fmt.Printf("[0m")

			case 1:
				fmt.Printf("\x1B")
				fmt.Printf("[35m")
				fmt.Print("0 ")
				fmt.Printf("\033")
				fmt.Printf("[0m")
			}
		}
		fmt.Println("")
	}

	// Print column headers below the board if
	// less than 10. (double digits mess with spacing
	if len(board) < 10 {

		// Prints bottom row of board with (=)
		for m := 0; m < len(board); m++ {
			fmt.Print("= ")
		}

		fmt.Println("")

		// Prints the column numbers below board
		for m := 0; m < len(board); m++ {
			fmt.Print(m, " ")
		}
	}
	fmt.Println("")

}

// Updates global variables to last token dropped location
func last_move(row int, column int) {
	global_row = row
	global_column = column
}

// Checks win conditions to determine if player who placed token won
func winner(length_to_win int, board [][]int) int {

	currentC := global_column
	currentR := global_row
	win := -1

	// Check if there is a tie. Returns 2
	win = check_tie(board)
	if win != -1 {
		return 2
	}

	// Check if row win
	win = check_row_win(currentR, currentC, length_to_win, board)
	if win != -1 {
		return win
	}

	// Check if column win
	win = check_column_win(currentR, currentC, length_to_win, board)
	if win != -1 {
		return win
	}

	// Check if a forward diaginal win
	win = check_forward_diaginal_win(currentR, currentC, length_to_win, board)
	if win != -1 {
		return win
	}

	// Check if a backward diaginal win
	win = check_backward_diaginal_win(currentR, currentC, length_to_win, board)
	if win != -1 {
		return win
	}

	// if no win
	return -1
}

// Initial setup of board that makes everything empty(-1)
func initialize_board(board [][]int) {

	for n := 0; n < len(board); n++ {
		for m := 0; m < len(board[0]); m++ {
			board[n][m] = -1
		}
	}

	print_board(board)
}

func place_token(player int, column int, board [][]int) int {

	// Check if column is valid for input
	if column < len(board[0]) {

		// We know the column number. So we loop threw the rows starting
		// at the bottom and look for an open spot
		for n := len(board) - 1; n >= 0; n-- {

			// If open spot
			if board[n][column] == -1 {

				// Insert player token
				if player == 0 {
					board[n][column] = 0
				} else {
					board[n][column] = 1
				}

				// Print the board and set last move
				print_board(board)
				last_move(n, column)

				return 0
			}
		}

		// Column is full
		fmt.Println("")
		fmt.Println("Sorry column full. Please try a different column")
		return 1

		// Column does not exist
	} else {
		fmt.Println("")
		fmt.Println("Sorry column does not exist. Please try a different column")
		return 1
	}
	return 1
}
