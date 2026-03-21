package tg

import (
	"fmt"
	"log/slog"

	chatdomain "github.com/braunkc/ai-bot-constructor/bot-container/internal/domain/chat"
	openrouter "github.com/braunkc/ai-bot-constructor/bot-container/internal/infra/open_router"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot          *tgbotapi.BotAPI
	systemPrompt string
	openRouter   *openrouter.OpenRouter
	log          *slog.Logger
}

func NewBot(token, systemPrompt string, openRouter *openrouter.OpenRouter, log *slog.Logger) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}
	bot.Debug = true

	return &TelegramBot{
		bot:          bot,
		systemPrompt: systemPrompt,
		openRouter:   openRouter,
		log:          log,
	}, nil
}

func (tb *TelegramBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	systemMessage := chatdomain.NewMessage("system", tb.systemPrompt)

	updates := tb.bot.GetUpdatesChan(u)
	for update := range updates {
		go func() {
			if update.Message != nil {
				if update.Message.IsCommand() {
					command := update.Message.Command()
					switch command {
					case "start":
						tb.SendMsg(&update, tb.systemPrompt)
					}

					return
				}

				tb.log.Debug("received message")
				waitMsg, err := tb.SendMsg(&update, "Секунду...")
				if err != nil {
					tb.log.Error("failed to send msg", slog.Any("err", err))
					return
				}

				messages := []chatdomain.Message{
					systemMessage,
					chatdomain.NewMessage("user", update.Message.Text),
				}

				respMsg, err := tb.openRouter.CreateResponse(messages)
				if err != nil {
					tb.log.Error("failed to create open router response", slog.Any("err", err))
					return
				}

				if err := tb.DeleteMsg(waitMsg); err != nil {
					tb.log.Error("failed to delete msg", slog.Any("err", err))
					return
				}

				if _, err := tb.SendMsg(&update, respMsg.Content); err != nil {
					tb.log.Error("failed to send telegram message", slog.Any("err", err))
					return
				}
			}
		}()
	}
}

func (tb *TelegramBot) SendMsg(update *tgbotapi.Update, text string) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID

	message, err := tb.bot.Send(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send msg: %w", err)
	}

	return &message, nil
}

func (tb *TelegramBot) DeleteMsg(msg *tgbotapi.Message) error {
	msgToDelete := tgbotapi.DeleteMessageConfig{
		ChatID:    msg.Chat.ID,
		MessageID: msg.MessageID,
	}

	if _, err := tb.bot.Request(msgToDelete); err != nil {
		return fmt.Errorf("failed to delete msg: %w", err)
	}

	return nil
}
