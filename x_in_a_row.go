package main

// Determines if the given slice (pointer to array 
// with self contained length method) contains a 
// winning value. 
func x_in_a_row(slice []int, win_length int) int {
    
    
    var count int = 0
    
    // loop threw the array checking for blocks of win length
    // equality 
    for i := 0; i < (len(slice) - (win_length - 1)); i++ {
	
	// No need to check empty spaces (-1)
	if slice[i] != -1 {
	    
	    for k := 1; k < win_length; k++ {
	    
		if slice[i] == slice[i + k] {
		    count++
		} else {
		    count = 0
		}
	    }
	
	    if count == win_length-1 {
		return slice[i]
	    }
	}
    }
	return -1
}

/* For Testing only!
func main() {
    a := []int{1,-1,1,1,0,0,0,0}
    
    fmt.Println(x_in_a_row(a,4))
}
*/

