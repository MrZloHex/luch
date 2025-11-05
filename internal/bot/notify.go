package bot

import (
	"encoding/json"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "log/slog"
)

type Notifier struct {
	notifiers  map[int64]struct{}
	notifyMu   sync.RWMutex
	notifyFile string
}

func (bot *Bot) setupNotifier() error {
	bot.not.notifyMu.Lock()
	defer bot.not.notifyMu.Unlock()

	bot.not.notifiers = make(map[int64]struct{})

	data, err := os.ReadFile(bot.not.notifyFile)
	if err != nil {
		log.Error("Failed to read notifiers", "err", err)
		return err
	}

	var ids []int64
	if err := json.Unmarshal(data, &ids); err != nil {
		log.Warn("Failed to unmarshal notify file", "err", err)
		return err
	}
	for _, id := range ids {
		bot.not.notifiers[id] = struct{}{}
	}

	return nil
}

func (bot *Bot) saveNotifiers() error {
	bot.not.notifyMu.RLock()
	defer bot.not.notifyMu.RUnlock()

	ids := make([]int64, 0, len(bot.not.notifiers))
	for id := range bot.not.notifiers {
		ids = append(ids, id)
	}

	data, err := json.Marshal(ids)
	if err != nil {
		log.Error("Failed to marshal notify file", "err", err)
		return err
	}
	if err := os.WriteFile(bot.not.notifyFile, data, 0644); err != nil {
		log.Error("Failed to write notify file", "err", err)
		return err
	}

	return nil
}

func (bot *Bot) toggleNotify(id int64) string {
	bot.not.notifyMu.Lock()
	defer bot.not.notifyMu.Unlock()

	if _, ok := bot.not.notifiers[id]; ok {
		delete(bot.not.notifiers, id)
		return "Notifications disabled"
	}
	bot.not.notifiers[id] = struct{}{}
	return "Notifications enabled"
}

func (bot *Bot) NotifyAll(text string) {
	bot.not.notifyMu.RLock()
	defer bot.not.notifyMu.RUnlock()

	for id := range bot.not.notifiers {
		msg := tgbotapi.NewMessage(id, text)
		if _, err := bot.api.Send(msg); err != nil {
			log.Error("Failed to send notify", "err", err)
		}
	}
}

func (bot *Bot) NotifyAllWithMarkup(text string, markup tgbotapi.InlineKeyboardMarkup) {
	bot.not.notifyMu.RLock()
	defer bot.not.notifyMu.RUnlock()

	for id := range bot.not.notifiers {
		msg := tgbotapi.NewMessage(id, text)
		msg.ReplyMarkup = markup
		if _, err := bot.api.Send(msg); err != nil {
			log.Error("Failed to send notify", "err", err)
		}
	}
}
