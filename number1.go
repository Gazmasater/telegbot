package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI("6515068455:AAF8-PKB0axfRWDC15kr1qRn8rOGl9GdDdk")
	if err != nil {
		log.Fatal(err)
	}

	// Включение отладочного режима
	bot.Debug = true

	// Логгирование успешной авторизации
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Настройка обновлений от бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// Слежение за состоянием пользователя
	userStates := make(map[int64]string)

	// Обработка обновлений
	for update := range updates {
		// Обработка текстовых сообщений
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Обработка команды "/start"
			if strings.HasPrefix(update.Message.Text, "/start") {
				sendWelcomeMessage(bot, update)
			}
		} else if update.CallbackQuery != nil {
			// Обработка коллбеков (нажатий на кнопки)
			handleCallbackQuery(bot, update, userStates)
		}
	}
}

// Отправка приветственного сообщения
func sendWelcomeMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	welcomeMessage := "Здравствуйте, Вас приветствует компания ГАЗмастер! В каком котле у вас возникла проблема?"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)

	// Создание клавиатуры с типами котлов
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Газовый", "gas"),
			tgbotapi.NewInlineKeyboardButtonData("Электрический", "electric"),
		),
	)
	msg.ReplyMarkup = keyboard

	// Отправка приветственного сообщения
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// Обработка коллбеков (нажатий на кнопки)
func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, userStates map[int64]string) {
	var replyMessage string

	switch update.CallbackQuery.Data {
	case "gas":
		if userStates[update.CallbackQuery.Message.Chat.ID] == "gas" {
			return
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

	case "floor", "wall", "electric":
		// Установите состояние пользователя для следующего шага
		userStates[update.CallbackQuery.Message.Chat.ID] = update.CallbackQuery.Data

		// Задайте вопрос о проблеме
		replyMessage = "Напишите, пожалуйста, какая у вас проблема с котлом?"

	default:
		replyMessage = "Неизвестная команда. Попробуйте снова."
	}

	if replyMessage != "" {
		// Отправка сообщения с результатом выбора
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replyMessage)
		bot.Send(msg)
	}

	// Подтверждение обработки коллбека
	callbackMsg := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	bot.AnswerCallbackQuery(callbackMsg)
}
