package main

import (
    "fmt"
    "os"
    "strconv"
)


func take_turns(length_to_win int, board[][]int) int {
    
	win := -1
	var column int
	error := 0
	
	for win == -1 {
	    
	    // Prompt player 1
	    fmt.Printf("\x1B")
	    fmt.Printf("[36m")
	    fmt.Print("Player 1 enter a column: ")
	    fmt.Printf("\033")
            fmt.Printf("[0m")
	    fmt.Scanf("%d", &column)
	    
	    
	    error = place_token(0, column, board)
	    
	    // Repeat if there was an error
	    for error == 1 {
		 fmt.Printf("\x1B")
	         fmt.Printf("[36m")
		 fmt.Print("Player 1 enter a column: ")
		 fmt.Printf("\033")
		 fmt.Printf("[0m")
		 fmt.Scanf("%d", &column)
		 error = place_token(0, column, board)
	    }
	    
	    // Checks if player won
	    win = winner(length_to_win,board)
	    
	    // If win
	    if win == 0 {
		return 0
	    
	    // If tie
	    } else if win == 2 {
		return 2
	    }
	    
	    
	    // Prompt player 2
	    fmt.Printf("\x1B")
	    fmt.Printf("[35m")
	    fmt.Print("Player 2 enter a column: ")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    fmt.Scanf("%d", &column)
	    
	    error = place_token(1, column, board)
	    
	    // Repeat if there was an error
	    for error == 1 {
		 fmt.Printf("\x1B")
	         fmt.Printf("[35m")
		 fmt.Print("Player 2 enter a column: ")
		 fmt.Printf("\033")
	         fmt.Printf("[0m")
		 fmt.Scanf("%d", &column)
		 error = place_token(1, column, board)
	    }
	    
	     // Checks if player won
	    win = winner(length_to_win,board)
	    
	    
	    // If win
	    if win == 0 {
		return 1
	    
	    // If tie
	    } else if win == 2 {
		return 2
	    }
	    
	    
	}
	return -1
}
// Main method 
func main() {
    
    // Get system arguements 
    args := os.Args;

    
    // If no commands then default board 
    // and length to win sizes
    if len(args) == 1 {
	
	// The gameboard is a slice, which is a reference to an array
	// So we will have a slice reference to another slice, similar
	// to the double array style in C. First we use the make function
	// to establish the total array. We will have a reference to an array of 
	// arrays in the end
	game := make([][]int,8)
    
	// Now to finish it we will loop threw the array and make the array at each
	// index 
	for i := 0; i < 8; i++ {
	    game[i] = make([]int,8)
	}
   
	length_to_win := 4
	win := -1
	
	initialize_board(game)
	
	for win == -1 {
	    
	    win = take_turns(length_to_win,game)
	}
	
	// If win
	if win == 0 {
	    fmt.Printf("\x1B")
	    fmt.Printf("[36m")
	    fmt.Println("Player 1 wins!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    return 
	} else if win == 1 {
	    fmt.Printf("\x1B")
	    fmt.Printf("[35m")
	    fmt.Println("Player 2 wins!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    
	    return 
	
	// If tie
	} else if win == 2 {
	    fmt.Println("Tie game!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    return 
	}
	    
    // Custom square board size and target to win
    } else if len(args) == 3 {
	
	size,_ := strconv.Atoi(args[1])
	game := make([][]int,size)
    
	// Now to finish it we will loop threw the array and make the array at each
	// index 
	for i := 0; i < size; i++ {
	    game[i] = make([]int,size)
	}
   
	length_to_win,_ := strconv.Atoi(args[2])
	win := -1
	
	initialize_board(game)
	
	for win == -1 {
	    
	    win = take_turns(length_to_win,game)
	}
	
	// If win
	if win == 0 {
	    fmt.Printf("\x1B")
	    fmt.Printf("[36m")
	    fmt.Println("Player 1 wins!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    return 
	} else if win == 1 {
	    fmt.Printf("\x1B")
	    fmt.Printf("[35m")
	    fmt.Println("Player 2 wins!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    
	    return 
	
	// If tie
	} else if win == 2 {
	    fmt.Println("Tie game!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    return 
	}
	
    } else if len(args) == 4 {
	
	rows, _ := strconv.Atoi(args[1])
	colum, _ := strconv.Atoi(args[2])
	game := make([][]int,rows)
    
	// Now to finish it we will loop threw the array and make the array at each
	// index 
	for i := 0; i < rows; i++ {
	    game[i] = make([]int,colum)
	}
   
	length_to_win, _ := strconv.Atoi(args[3])
	win := -1
	
	initialize_board(game)
	
	for win == -1 {
	    
	    win = take_turns(length_to_win,game)
	}
	
	// If win
	if win == 0 {
	    fmt.Printf("\x1B")
	    fmt.Printf("[36m")
	    fmt.Println("Player 1 wins!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    return 
	} else if win == 1 {
	    fmt.Printf("\x1B")
	    fmt.Printf("[35m")
	    fmt.Println("Player 2 wins!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    
	    return 
	
	// If tie
	} else if win == 2 {
	    fmt.Println("Tie game!")
	    fmt.Printf("\033")
	    fmt.Printf("[0m")
	    return 
	}
	
    // Invalid arguments
    } else {
	fmt.Println("Sorry invalid commands")
	return 
    }
    return
}