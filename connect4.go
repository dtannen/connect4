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
	"strings"
)

func take_turns(length_to_win int, board [][]int, lst *list.List, first_player int) int {

	win := -1
	error := 0
	var player1 ClientPlayer
	var player2 ClientPlayer
	for e := lst.Front(); e != nil; e = e.Next() {
		client := e.Value.(ClientPlayer)
		if client.Name == "Player1" {
			if first_player == 1 {
				player1 = client
			} else {
				player2 = client
			}
		}
		if client.Name == "Player2" {
			if first_player == 2 {
				player1 = client
			} else {
				player2 = client
			}
		}
		// io.Copy(*client.Con, bytes.NewBufferString(message))
	}

	for win == -1 {

		// Prompt player 1
		fmt.Printf("\x1B")
		fmt.Printf("[36m")
		fmt.Print("Player 1 enter a column: ")
		fmt.Printf("\033")
		fmt.Printf("[0m")
		// get column number from first_player channel

		column_string := <-player1.OUT
		if column_string == "quit" {
			return 3
			break
		}
		column, err := strconv.Atoi(strings.TrimSpace(column_string))
		if err != nil {
			column = 9
			io.Copy(*player1.Con, bytes.NewBufferString("err\n"))
		}
		// fmt.Scanf("%d", &column)

		error = place_token(0, column, board)

		// Repeat if there was an error
		for error == 1 {
			fmt.Printf("\x1B")
			fmt.Printf("[36m")
			fmt.Print("Player 1 enter a column: ")
			fmt.Printf("\033")
			fmt.Printf("[0m")
			column_string := <-player1.OUT
			if column_string == "quit" {
				return 3
				break
			}
			column, err := strconv.Atoi(strings.TrimSpace(column_string))
			if err != nil {
				column = 9
				io.Copy(*player1.Con, bytes.NewBufferString("err\n"))
			}
			// fmt.Scanf("%d", &column)
			error = place_token(0, column, board)
		}

		// inform player2
		io.Copy(*player2.Con, bytes.NewBufferString(column_string))

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
		column_string = <-player2.OUT
		if column_string == "quit" {
			return 3
			break
		}
		column, err = strconv.Atoi(strings.TrimSpace(column_string))
		if err != nil {
			column = 9
			io.Copy(*player2.Con, bytes.NewBufferString("err\n"))
		}
		// fmt.Scanf("%d", &column)

		error = place_token(1, column, board)

		// Repeat if there was an error
		for error == 1 {
			fmt.Printf("\x1B")
			fmt.Printf("[35m")
			fmt.Print("Player 2 enter a column: ")
			fmt.Printf("\033")
			fmt.Printf("[0m")
			column_string := <-player2.OUT
			if column_string == "quit" {
				return 3
				break
			}
			column, err := strconv.Atoi(strings.TrimSpace(column_string))
			if err != nil {
				column = 9
				io.Copy(*player2.Con, bytes.NewBufferString("err\n"))
			}
			// fmt.Scanf("%d", &column)
			error = place_token(1, column, board)
		}

		// inform player1
		io.Copy(*player1.Con, bytes.NewBufferString(column_string))

		// Checks if player won
		win = winner(length_to_win, board)

		// If win
		if win == 1 {
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
	for e := c.ListChain.Back(); e != nil; e = e.Prev() {
		client := e.Value.(ClientPlayer)
		fmt.Println(client.Name)
		(*client.Con).Close()
		client.ListChain.Remove(e)
		client.IN <- "quit"
		client.OUT <- "quit"
		close(client.OUT)
	}
}

func request_handler(conn *net.Conn, out chan string, lst *list.List, connections *int) {
	channel := make(chan string)
	out_channel := make(chan string)
	// defer close(channel)
	// add listener for channel to send msgs to player
	go send_player_data(channel, conn)
	playername := "Player" + strconv.Itoa(lst.Len()+1)
	newclient := &ClientPlayer{playername, conn, channel, out_channel, lst}
	channel <- playername + "\n"
	lst.PushBack(*newclient)
	fmt.Println(lst.Len())
	if lst.Len() == 2 {
		go start_game(1, out, lst)
	}
	for {
		msg, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			newclient.Close()
			*connections = 0
			break
		}
		if msg == "quit\r\n" || msg == "quit\n" {
			fmt.Println("end game.")
			newclient.Close()
			*connections = 0
			break
		}
		out_channel <- string(msg) + "\n"
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

func start_game(first_player int, out chan<- string, lst *list.List) {
	// initialize game board
	game := make([][]int, 8)
	for i := 0; i < 8; i++ {
		game[i] = make([]int, 8)
	}
	length_to_win := 4
	win := -1
	initialize_board(game)
	for win == -1 {
		win = take_turns(length_to_win, game, lst, first_player)
	}
	// If win
	if win == 0 {
		// out <- "Player1 wins!\n"
		fmt.Printf("\x1B")
		fmt.Printf("[36m")
		fmt.Println("Player 1 wins!")
		fmt.Printf("\033")
		fmt.Printf("[0m")
	} else if win == 1 {
		// out <- "Player2 wins!\n"
		fmt.Printf("\x1B")
		fmt.Printf("[35m")
		fmt.Println("Player 2 wins!")
		fmt.Printf("\033")
		fmt.Printf("[0m")
		// If tie
	} else if win == 2 {
		// out <- "Tie game!\n"
		fmt.Println("Tie game!")
		fmt.Printf("\033")
		fmt.Printf("[0m")
	}
	if win == 3 {
		fmt.Println("closing game")
	} else {
		if first_player == 1 {
			start_game(2, out, lst)
		} else {
			start_game(1, out, lst)
		}
	}
}

// Main method
func main() {

	// Get system arguements
	args := os.Args
	connections := 0
	// If no commands then default board
	// and length to win sizes
	if len(args) == 1 {
		psock, err := net.Listen("tcp", ":3000")
		if err != nil {
			// handle error
			fmt.Println("Can't start server!")
		}
		clientlist := list.New()
		channel := make(chan string)
		go send_data(channel, clientlist)
		for {
			if connections < 2 {
				connections++
				conn, err := psock.Accept()
				if err != nil {
					return
				}
				go request_handler(&conn, channel, clientlist, &connections)
			} else {
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
