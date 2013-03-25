package main



// Checks for row wins (-). Returns -1 for no win,
// 0 for player 1 win, 1 for player 2 win
func check_row_win(currentR int, currentC int, length_to_win int, board[][]int) int {
    
    // Temp row array to check
    check_row := make([]int,len(board[0]))
    
    // Create the row from the last dropped token
    for m := 0; m < len(board[0]); m++ {
	check_row[m] = board[currentR][m]
    }
    
    return x_in_a_row(check_row, length_to_win)
}

// Checks for column wins (|). Returns -1 for no win,
// 0 for player 1 win, 1 for player 2 wins
func check_column_win(currentR int, currentC int, length_to_win int, board[][]int) int {
    
    // Temp column array to check
    check_column := make([]int,len(board))
    
    // Create the column from the last dropped token
    for n := 0; n < len(board); n++ {
	check_column[n] = board[n][currentC]
    }
    
    return x_in_a_row(check_column,length_to_win)
}

// Checks for forward diaginal win (\). Returns -1 for no win,
// 0 for player 1 win, 1 for player 2 win
func check_forward_diaginal_win(currentR int, currentC int, length_to_win int, board[][]int) int {
    
    // Temp array to check win
    check_forw_diaginal := make([]int,len(board[0]))
    size := 0
    
    // Check bottom half of board for left diaginal (\)
    if currentR > currentC {
	
	// Get the row number of the first diaginal 
	for currentC != 0 {
	    currentR--
	    currentC--
	}
	
	// Start from the spot of the first diaginal. 
	// Then decrement down the diaginal bottom
	for m:= 0; currentR < len(board); m++  {
	    check_forw_diaginal[m] = board[currentR][m]
	    currentR++
	    size++
	}
	
    // Check to see if it is on the upper part of the board
    } else if currentC > currentR {
	
	for currentR != 0 {
	    currentR--
	    currentC--
	}
	
	for m:= 0; currentC < len(board[0]); m++  {
	    check_forw_diaginal[m] = board[m][currentC]
	    currentC++
	    size++
	}
	
    // We know that it falls in the middle diaginal
    } else {
	
	for m := 0; m < len(board[0]); m++ {
	    check_forw_diaginal[m] = board[m][m]
	    size++;
	}
    }
    
    // Only need to check the adjusted size of the slice
    check_forw_diaginal = check_forw_diaginal[:size]
    
    return x_in_a_row(check_forw_diaginal,length_to_win)
}
	
// Checks for backward diaginal win (/). Returns -1 for no win,
// 0 for player 1 win, 1 for player 2 win
func check_backward_diaginal_win(currentR int, currentC int, length_to_win int, board[][]int) int {
    
    // Temp array to check for wins
    check_back_diaginal := make([]int,len(board[0]))
    
    // Size of slice
    size := 0
    
    // While the left and bottom board boundaries are not hit. If 
    // initial spot is in column 0 then move on. If initial spot
    // is not in column zero then try and find the bottom row
    // or until column zero is hit.
    for currentC > 0 && currentR < len(board)-1 {
	currentC--
	currentR++
    }
   
    // Gets the diaginal from the diaginal spot starting at the bottom left.
    // We keep getting the diaginal position and adding it to a temp array
    // until either the top row is hit or the farthest right column is hit
    for m:= 0; currentC < len(board[0]) && currentR >= 0; m++ {
	check_back_diaginal[m] = board[currentR][currentC]
	currentC++
	currentR-- 
	size++
    }
    
    // We only need to check the slice from 0 to size. The remaining would 
    // be initialized 0s that would be off of the board
    check_back_diaginal = check_back_diaginal[:size]
 
    return x_in_a_row(check_back_diaginal,length_to_win)
}

// Checks if the game board is full and there is a tie.
// Returns -1 for no tie and 2 for tie
func check_tie(board[][]int) int {
    
    var tie = 2
    
    for n := 0; n < len(board); n++ {
	for m := 0; m < len(board[0]); m++ {
	    
	    // If there is an open spot left there is no tie
	    if board[n][m] == -1 {
		tie = -1
	    }
	}
    }
    return tie
}
    