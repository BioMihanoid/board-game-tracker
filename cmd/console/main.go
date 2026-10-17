package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	consolehandler "board-game-tracker/internal/handler/console"
	mapdb "board-game-tracker/internal/repository/map"
	"board-game-tracker/internal/service"
)

func main() {
	ctx := context.Background()

	gameRepo := mapdb.NewGameRepository()
	gameService := service.NewGameService(gameRepo)

	handler := consolehandler.NewHandler(gameService)

	choice := bufio.NewScanner(os.Stdin)

	fmt.Println("Board Game Tracker")
	for {
		fmt.Println("1. Добавить игру")
		fmt.Println("2. Показать все игры")
		fmt.Println("0. Выход")

		fmt.Print("Введите номер действия: ")

		choice.Scan()
		title := choice.Text()

		switch title {
		case "1":
			err := handler.AddGame(ctx)
			if err != nil {
				fmt.Println("Не получилось добавить игру")
			}
		case "2":
			handler.PrintAllGames(ctx)
		case "0":
			return
		default:
			fmt.Printf("Ошибка, введите корректную цифру!\n")
		}
	}
}
