package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func take_turns(length_to_win int, board [][]int) int {

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
		win = winner(length_to_win, board)

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
		win = winner(length_to_win, board)

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

// compare two clients: name and network connection
func (c *ClientPlayer) Equal(cl *ClientPlayer) bool {
	if c.Con == cl.Con {
		return true
	}
	return false
}

// delete the client from list
func (c *ClientPlayer) Close() {
	for e := c.ListChain.Front(); e != nil; e = e.Next() {
		client := e.Value.(ClientPlayer)
		(*client.Con).Close()
		client.ListChain.Remove(e)
	}
}

func request_handler(conn *net.Conn, out chan string, lst *list.List) {
	// 	defer close(out)
	newclient := &ClientPlayer{conn, lst}
	lst.PushBack(*newclient)
	for {
		msg, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			newclient.Close()
			break
		}
		if msg == "quit\r\n" || msg == "quit\n" {
			fmt.Println("quitting...")
			newclient.Close()
			break
		}
		out <- string(msg) + "\n"
	}
}

func send_data(in <-chan string, lst *list.List) {
	for {
		message := <-in
		if lst.Len() > 0 {
			log.Print(message)
			for e := lst.Front(); e != nil; e = e.Next() {
				client := e.Value.(ClientPlayer)
				io.Copy(*client.Con, bytes.NewBufferString(message))
			}
		}
	}
}

type ClientPlayer struct {
	Con       *net.Conn  // connection of client
	ListChain *list.List // reference to list
}

// Main method
func main() {

	// Get system arguements
	args := os.Args

	// If no commands then default board
	// and length to win sizes
	if len(args) == 1 {

		// The gameboard is a slice, which is a reference to an array
		// So we will have a slice reference to another slice, similar
		// to the double array style in C. First we use the make function
		// to establish the total array. We will have a reference to an array of
		// arrays in the end
		game := make([][]int, 8)

		// Now to finish it we will loop threw the array and make the array at each
		// index
		for i := 0; i < 8; i++ {
			game[i] = make([]int, 8)
		}

		length_to_win := 4
		win := -1

		initialize_board(game)

		psock, err := net.Listen("tcp", ":3000")
		if err != nil {
			// handle error
			fmt.Println("Can't start server!")
		}
		clientlist := list.New()
		channel := make(chan string)
		go send_data(channel, clientlist)
		for {
			if clientlist.Len() < 1 {
				conn, err := psock.Accept()
				if err != nil {
					return
				}
				go request_handler(&conn, channel, clientlist)
			} else {
				conn, err := psock.Accept()
				if err != nil {
					return
				}
				io.Copy(conn, bytes.NewBufferString("Game In Progress\n"))
				conn.Close()
			}
		}

		for win == -1 {
			win = take_turns(length_to_win, game)
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
	} else {
		fmt.Println("Sorry invalid commands")
		return
	}
	return
}
