package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	consolehandler "board-game-tracker/internal/handler/console"
	memoryrepo "board-game-tracker/internal/repository/memory"
	"board-game-tracker/internal/service"
)

func main() {
	ctx := context.Background()

	gameRepo := memoryrepo.NewGameRepository()
	gameService := service.NewGameService(gameRepo)

	scanner := bufio.NewScanner(os.Stdin)
	handler := consolehandler.NewHandler(gameService, scanner)

	for {
		printMenu()

		if !scanner.Scan() {
			fmt.Printf("\nВвод завершён, выходим\n")
			return
		}

		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			id, err := handler.AddGame(ctx)
			if err != nil {
				if quitOnError(err) {
					return
				}

				continue
			}

			fmt.Printf("Игра успешно добавлена с id: %d\n\n", id)
		case "2":
			fmt.Println()
			if err := handler.PrintBaseGames(ctx); err != nil {
				if quitOnError(err) {
					return
				}
			}
		case "3":
			fmt.Println()
			if err := handler.PrintExpansions(ctx); err != nil {
				if quitOnError(err) {
					return
				}
			}
		case "4":
			fmt.Println()
			if err := handler.PrintGameByID(ctx); err != nil {
				if quitOnError(err) {
					return
				}
			}
		case "0":
			return
		default:
			fmt.Printf("Ошибка, введите корректную цифру!\n\n")
		}
	}
}

func printMenu() {
	fmt.Printf("Board Game Tracker\n\n")
	fmt.Printf("1. Добавить игру\n")
	fmt.Printf("2. Показать все игры в коллекции\n")
	fmt.Printf("3. Показать все дополнения по ID игры\n")
	fmt.Printf("4. Показать игру по ID\n")
	fmt.Printf("0. Выход\n\n")
	fmt.Printf("Введите номер действия: ")
}

// quitOnError печатает ошибку и сообщает, пора ли завершать программу.
func quitOnError(err error) bool {
	if errors.Is(err, consolehandler.ErrInputClosed) {
		fmt.Printf("\nВвод завершён, выходим\n")
		return true
	}

	fmt.Printf("Ошибка: %v\n\n", err)

	return false
}
