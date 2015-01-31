package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
		client.IN <- "quit"
		close(client.OUT)
	}
}

func request_handler(conn *net.Conn, out chan string, lst *list.List) {
	channel := make(chan string)
	out_channel := make(chan string)
	// defer close(channel)
	// add listener for channel to send msgs to player
	go send_player_data(channel, conn)
	playername := "Player" + strconv.Itoa(lst.Len()+1)
	newclient := &ClientPlayer{playername, conn, channel, out_channel, lst}
	channel <- playername + "\n"
	lst.PushBack(*newclient)
	for {
		msg, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			newclient.Close()
			break
		}
		if msg == "quit\r\n" || msg == "quit\n" {
			fmt.Println("end game.")
			newclient.Close()
			break
		}
		out <- string(msg) + "\n"
	}
}

func send_player_data(in chan string, conn *net.Conn) {
	for {
		message := <-in
		if message == "quit" {
			close(in)
			break
		}
		io.Copy(*conn, bytes.NewBufferString(message))
	}
}

func send_data(in <-chan string, lst *list.List) {
	for {
		message := <-in
		if lst.Len() > 0 {
			for e := lst.Front(); e != nil; e = e.Next() {
				client := e.Value.(ClientPlayer)
				io.Copy(*client.Con, bytes.NewBufferString(message))
			}
		}
	}
}

type ClientPlayer struct {
	Name      string      // players name
	Con       *net.Conn   // connection of client
	IN        chan string // channel to send messages to user
	OUT       chan string // channel to get messages from user
	ListChain *list.List  // reference to list
}

func start_game(out chan<- string, lst *list.List) {
	// initialize game board
	game := make([][]int, 8)
	for i := 0; i < 8; i++ {
		game[i] = make([]int, 8)
	}
	length_to_win := 4
	win := -1
	initialize_board(game)
	for win == -1 {
		win = take_turns(length_to_win, game)
	}
	// If win
	if win == 0 {
		out <- "Player 1 wins!\n"
		fmt.Printf("\x1B")
		fmt.Printf("[36m")
		fmt.Println("Player 1 wins!")
		fmt.Printf("\033")
		fmt.Printf("[0m")
	} else if win == 1 {
		out <- "Player 2 wins!\n"
		fmt.Printf("\x1B")
		fmt.Printf("[35m")
		fmt.Println("Player 2 wins!")
		fmt.Printf("\033")
		fmt.Printf("[0m")
		// If tie
	} else if win == 2 {
		out <- "Tie game!\n"
		fmt.Println("Tie game!")
		fmt.Printf("\033")
		fmt.Printf("[0m")
	}
	start_game(out, lst)
}

// Main method
func main() {

	// Get system arguements
	args := os.Args

	// If no commands then default board
	// and length to win sizes
	if len(args) == 1 {
		psock, err := net.Listen("tcp", ":3000")
		if err != nil {
			// handle error
			fmt.Println("Can't start server!")
		}
		clientlist := list.New()
		game_started := -1
		channel := make(chan string)
		go send_data(channel, clientlist)
		for {
			if clientlist.Len() < 1 {
				game_started = -1
				conn, err := psock.Accept()
				if err != nil {
					return
				}
				go request_handler(&conn, channel, clientlist)
			} else {
				if game_started == -1 {
					go start_game(channel, clientlist)
					game_started = 1
				}
				conn, err := psock.Accept()
				if err != nil {
					return
				}
				io.Copy(conn, bytes.NewBufferString("Game In Progress\n"))
				conn.Close()
			}
		}

	} else {
		fmt.Println("Sorry invalid commands")
		return
	}
	return
}
