package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	mapdb "board-game-tracker/internal/repository/map"
	"board-game-tracker/internal/service"
)

func main() {
	gameRepo := mapdb.NewGameRepository()
	gameService := service.NewGameService(gameRepo)

	fmt.Println("Board Game Tracker")
	for {
		fmt.Println("1. Добавить игру")
		fmt.Println("2. Показать все игры")
		fmt.Println("0. Выход")

		choice := bufio.NewScanner(os.Stdin)

		fmt.Print("Введите номер действия: ")
		choice.Scan()
		title := choice.Text()

		switch title {
		case "1":
			fmt.Println("Добавили игру...")
		case "2":
			gameService.PrintGames(context.Background())
		case "0":
			return
		}
	}
}
