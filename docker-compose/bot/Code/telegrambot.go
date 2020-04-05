package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"os"
	"reflect"
	"time"
)

func telegramBot() {
	START_MESSAGE := "Hi, i'm a wikipedia bot, i can search information in a wikipedia, send me something what you want find in Wikipedia."
	DB_ERROR_MESSAGE := "Database error."
	USER_COUNT_MESSAGE := "%d peoples used me for search information in Wikipedia"
	DB_NOT_CONNECTED_MESSAGE := "Database not connected, so i can't say you how many peoples used me."

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			chatId := update.Message.Chat.ID

			switch update.Message.Text {

			case "/start":
				msg := tgbotapi.NewMessage(chatId, START_MESSAGE)
				bot.Send(msg)
			case "/number_of_users":
				if os.Getenv("DB_SWITCH") == "on" {
					num, err := getNumberOfUsers()
					if err != nil {
						msg := tgbotapi.NewMessage(chatId, DB_ERROR_MESSAGE)
						bot.Send(msg)
					}

					ans := fmt.Sprint(USER_COUNT_MESSAGE, num)
					msg := tgbotapi.NewMessage(chatId, ans)
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(chatId, DB_NOT_CONNECTED_MESSAGE)
					bot.Send(msg)
				}
			default:
				language := os.Getenv("LANGUAGE")
				ms, _ := urlEncoded(update.Message.Text)
				url := ms
				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"
				message := wikipediaAPI(request)

				if os.Getenv("DB_SWITCH") == "on" {
					err := collectData(update.Message.Chat.UserName, chatId, update.Message.Text, message)
					if err != nil {
						msg := tgbotapi.NewMessage(chatId, DB_ERROR_MESSAGE)
						bot.Send(msg)
					}
				}
				for _, val := range message {

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(chatId, val)
					bot.Send(msg)
				}
			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}

func main() {

	time.Sleep(1 * time.Minute)

	//Создаем таблицу
	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := createTable(); err != nil {

				panic(err)
			}
		}
	}

	time.Sleep(1 * time.Minute)

	//Вызываем бота
	telegramBot()
}
