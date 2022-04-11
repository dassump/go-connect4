package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	name         = "🔵🔵🔵🔵 Connect4 🔴🔴🔴🔴"
	objective    = "💡 Objective: Be the first player to get four of your colored checkers in a row horizontally, vertically or diagonally."
	rules        = "📋 Rules: Try to build a row of four checkers while keeping your opponent from doing the same."
	help         = "🚀 How to play: Type the column number '1..7' to play or 'e' to exit and press <enter>."
	players      = "😃 Players: 2-Player"
	round        = "🥊 Round:"
	turn         = "🤩 Turn:"
	input_column = "✨ Enter your column: "
	input_quit   = "e"
	error_char   = "❗ Invalid character"
	error_column = "❗ Invalid column"
	end_draw     = "👍 The game ended in a draw, congratulations to all players."
	end_winner   = "🎉 Congratulations %s %s, you win the game!!!\n"
	end_bye      = "👋 Bye-Bye"
)

type (
	player struct {
		name   string
		avatar string
	}

	game struct {
		board     [7][6]*player
		winner    *player
		turn      *player
		round     int
		round_max int
	}
)

var (
	player1 = &player{"Player 1", "🔵"}
	player2 = &player{"Player 2", "🔴"}
)

func (game *game) Start() {
	game.turn = player1
	game.round = 1
	game.round_max = len(game.board) * len(game.board[0])
}

func (game *game) ChangeTurn() {
	switch game.turn {
	case player1:
		game.turn = player2
	case player2:
		game.turn = player1
	}

	game.round++
}

func (game *game) GetPlayer(column, line int) *player {
	return game.board[column][line]
}

func (game *game) ColumnFill(column int) (fill int) {
	for _, player := range game.board[column] {
		if player != nil {
			fill++
		}
	}

	return
}

func (game *game) LineAvatars(line int) (avatars []interface{}) {
	for column := 0; column < len(game.board); column++ {
		player := game.GetPlayer(column, line)
		if player != nil {
			avatars = append(avatars, player.avatar)
			continue
		}
		avatars = append(avatars, "  ")
	}

	return
}

func (game *game) Play(column int) {
	for line, player := range game.board[column-1] {
		if player == nil {
			game.board[column-1][line] = game.turn
			break
		}
	}
}

func (game *game) Draw() bool {
	return (game.round > game.round_max) && game.winner == nil
}

func (game *game) Winner() bool {
	for column := 0; column < len(game.board); column++ {
		for line := 0; line < len(game.board[0])-3; line++ {
			player := game.GetPlayer(column, line)
			if player != nil &&
				player == game.GetPlayer(column, line+1) &&
				player == game.GetPlayer(column, line+2) &&
				player == game.GetPlayer(column, line+3) {
				game.winner = player
				break
			}
		}

	}

	for line := 0; line < len(game.board[0]); line++ {
		for column := 0; column < len(game.board)-3; column++ {
			player := game.GetPlayer(column, line)
			if player != nil &&
				player == game.GetPlayer(column+1, line) &&
				player == game.GetPlayer(column+2, line) &&
				player == game.GetPlayer(column+3, line) {
				game.winner = player
				break
			}
		}
	}

	for column := 0; column < len(game.board)-3; column++ {
		for line := 0; line < len(game.board[0])-3; line++ {
			player := game.GetPlayer(column, line)
			if player != nil &&
				player == game.GetPlayer(column+1, line+1) &&
				player == game.GetPlayer(column+2, line+2) &&
				player == game.GetPlayer(column+3, line+3) {
				game.winner = player
				break
			}
		}

		for line := 3; line < len(game.board[0]); line++ {
			player := game.GetPlayer(column, line)
			if player != nil &&
				player == game.GetPlayer(column+1, line-1) &&
				player == game.GetPlayer(column+2, line-2) &&
				player == game.GetPlayer(column+3, line-3) {
				game.winner = player
				break
			}
		}
	}

	return game.winner != nil
}

func (game *game) ShowBoard() {
	fmt.Print("\033[H\033[2J")

	fmt.Printf("┏%v┓\n", strings.Repeat("━", 28))
	fmt.Printf("┃ %s ┃\n", name)
	fmt.Printf("┗%v┛\n", strings.Repeat("━", 28))

	fmt.Printf("\n%s\n%s\n%s\n%s\n", objective, rules, players, help)
	fmt.Printf("\n%s %d\n%s %s %s\n", round, game.round, turn, game.turn.avatar, game.turn.name)

	fmt.Print("\n    ")
	for column := 0; column <= len(game.board)-1; column++ {
		fmt.Print(column+1, "\ufe0f\u20e3", "    ")
	}
	fmt.Printf("\n  ┏━━%v━━┓\n", strings.Repeat("━━┳━━", 6))
	for line := len(game.board[0]) - 1; line >= 0; line-- {
		fmt.Print(line+1, "\ufe0f\u20e3", " ")
		fmt.Printf("┃ %s ┃ %s ┃ %s ┃ %s ┃ %s ┃ %s ┃ %s ┃\n", game.LineAvatars(line)...)
		if line > 0 {
			fmt.Printf("  ┣━━%v━━┫\n", strings.Repeat("━━╋━━", 6))
		}
	}
	fmt.Printf("  ┗━━%v━━┛\n\n", strings.Repeat("━━┻━━", 6))
}

func (game *game) Input() (choice int, err error) {
	fmt.Print(input_column)

	stdin := bufio.NewReader(os.Stdin)
	_, err = fmt.Fscanln(stdin, &choice)
	if err != nil {
		char, _, _ := stdin.ReadRune()
		if strings.EqualFold(string(char), input_quit) {
			game.End()
			os.Exit(0)
		}

		err = errors.New(error_char)
		stdin.Discard(stdin.Buffered())

		return
	}

	if choice < 1 || choice > len(game.board) || game.ColumnFill(choice-1) == len(game.board[choice]) {
		err = errors.New(error_column)
	}

	return
}

func (game *game) End() {
	switch {
	case game.winner != nil:
		fmt.Printf(end_winner, game.winner.avatar, game.winner.name)
	case game.round > game.round_max:
		fmt.Println(end_draw)
	}

	fmt.Println(end_bye)
}

func main() {
	game := new(game)
	game.Start()

	for !game.Winner() && !game.Draw() {
		game.ShowBoard()

		choice, err := game.Input()
		if err != nil {
			fmt.Print(err)
			time.Sleep(time.Second)
			continue
		}

		game.Play(choice)
		game.ChangeTurn()
	}

	game.ShowBoard()
	game.End()
}
