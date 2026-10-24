package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	consolehandler "board-game-tracker/internal/handler/console"
	mapdb "board-game-tracker/internal/repository/memory"
	"board-game-tracker/internal/service"
)

func main() {
	ctx := context.Background()

	gameRepo := mapdb.NewGameRepository()

	gameService := service.NewGameService(gameRepo)

	scanner := bufio.NewScanner(os.Stdin)

	handler := consolehandler.NewHandler(gameService, scanner)

	for {
		fmt.Printf("Board Game Tracker\n\n")

		fmt.Printf("1. Добавить игру\n")
		fmt.Printf("2. Показать все игры в коллекции\n")
		fmt.Printf("3. Показать все дополнения по ID игры\n")
		fmt.Printf("4. Показать игру по ID\n")
		fmt.Printf("0. Выход\n\n")

		fmt.Printf("Введите номер действия: ")

		ok := scanner.Scan()
		if !ok {
			break
		}

		title := scanner.Text()

		switch title {
		case "1":
			id, err := handler.AddGame(ctx)
			if err != nil {
				fmt.Printf("Не получилось добавить игру, попробуйте ещё раз\n\n")
			}
			fmt.Printf("Игра успешно добавлена с id: %d\n\n", id)
		case "2":
			fmt.Println()
			handler.PrintAllGames(ctx)
			fmt.Println()
		case "3":
			fmt.Println()
			handler.PrintAllAdditionalContentByGameID(ctx)
			fmt.Println()
		case "4":
			fmt.Println()
			handler.PrintGameByID(ctx)
			fmt.Println()
		case "0":
			return
		default:
			fmt.Printf("Ошибка, введите корректную цифру!\n\n")
		}
	}
}
