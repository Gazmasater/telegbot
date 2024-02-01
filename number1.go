package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6515068455:AAF8-PKB0axfRWDC15kr1qRn8rOGl9GdDdk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	userStates := make(map[int64]string)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if strings.HasPrefix(update.Message.Text, "/start") {
				welcomeMessage := "Здравствуйте, Вас приветствует компания ГАЗмастер! В каком котле у вас возникла проблема?"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)

				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Газовый", "gas"),
						tgbotapi.NewInlineKeyboardButtonData("Электрический", "electric"),
					),
				)
				msg.ReplyMarkup = keyboard

				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		} else if update.CallbackQuery != nil {
			var replyMessage string

			switch update.CallbackQuery.Data {
			case "gas":
				if userStates[update.CallbackQuery.Message.Chat.ID] == "gas" {
					continue
				}
				userStates[update.CallbackQuery.Message.Chat.ID] = "gas"

				// Изменение клавиатуры в уже существующем сообщении
				gasKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Напольный", "floor"),
						tgbotapi.NewInlineKeyboardButtonData("Настенный", "wall"),
					),
				)
				editMsg := tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      update.CallbackQuery.Message.Chat.ID,
						MessageID:   update.CallbackQuery.Message.MessageID,
						ReplyMarkup: &gasKeyboard,
					},
					Text: "Выберите тип газового котла:",
				}
				bot.Send(editMsg)

			case "floor":
				replyMessage = "Спасибо за информацию. Мы начинаем работу над решением проблемы в напольном газовом котле."

			case "wall":
				replyMessage = "Спасибо за информацию. Мы начинаем работу над решением проблемы в настенном газовом котле."

			case "electric":
				replyMessage = "Спасибо за информацию. Мы начинаем работу над решением проблемы в электрическом котле."

			default:
				replyMessage = "Неизвестная команда. Попробуйте снова."
			}

			if replyMessage != "" {
				// Отправка сообщения с результатом выбора
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replyMessage)
				bot.Send(msg)
			}

			callbackMsg := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			bot.AnswerCallbackQuery(callbackMsg)
		}
	}
}
