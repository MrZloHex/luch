package bot

import (
	"fmt"
	stdlog "log"
	log "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"luch/pkg/protocol"
	_ "luch/pkg/audio"

	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type BotConfig struct {
	Token  string
	Debug  bool
	Logger *stdlog.Logger
	Notify string
}

type Bot struct {
	api *tgbotapi.BotAPI

	ptcl *protocol.Protocol

	cmds Commands
	kb   Keyboard
	not  Notifier
}

func NewBot(cfg BotConfig, ptcl *protocol.Protocol) (*Bot, error) {
	log.Debug("init telebot")

	bot := Bot{
		ptcl: ptcl,
		not: Notifier{
			notifyFile: cfg.Notify,
		},
	}

	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Error("Failed to init api", "err", err)
		return nil, err
	}
	tgbotapi.SetLogger(cfg.Logger)
	api.Debug = cfg.Debug

	bot.api = api

	log.Debug("authorised telebot")

	return &bot, nil
}

func (bot *Bot) Setup() {
	bot.setupNotifier()

	err := bot.fetchCommands()
	if err != nil {
		log.Error("Failed to retrive commads", "err", err)
	}
	log.Debug("Commands from Telegram", "cmd", bot.cmds)

	bot.setupKeyboard()
}

func (bot *Bot) SendReq(to, pay string) string {
	resp, err := bot.ptcl.Send(to, pay)
	if err != nil {
		return fmt.Sprintf("Failed to send request: %s", err.Error)
	} else {
		return string(resp)
	}
}

func downloadFile(url string, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func convertToWav(src string, dst string) error {
	// 16 kHz mono PCM is fine for Whisper
	cmd := exec.Command("ffmpeg", "-y", "-i", src, "-ac", "1", "-ar", "16000", dst)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg: %v\n%s", err, string(out))
	}
	return nil
}

func runWhisper(ctx context.Context, whisperBin string, modelPath string, wavPath string, lang string, translate bool) (string, error) {
	args := []string{"-m", modelPath, "-f", wavPath, "-otxt"}

	// auto language detection unless explicitly provided
	if lang == "" {
		args = append(args, "-l", "auto")
	} else {
		args = append(args, "-l", lang)
	}

	if translate {
		// translate input → English
		args = append(args, "-tr")
	}

	cmd := exec.CommandContext(ctx, whisperBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper: %v\n%s", err, string(out))
	}

	b, readErr := os.ReadFile(wavPath + ".txt")
	if readErr == nil {
		return strings.TrimSpace(string(b)), nil
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return "", errors.New("no transcription produced")
	}
	return s, nil
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Debug("Got smth", "from", update.Message.From.UserName, "text", update.Message.Text)

		switch {
		case update.Message.IsCommand():
			bot.processCmd(update)
			continue
		case bot.isKeyboard(update):
			bot.proccessKeyboard(update)
			continue
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No such thingy, sorry\nIf you implement it or contact developer\nSee /help")
			bot.api.Send(msg)
		}

		var fileID string
		if update.Message.Voice != nil {
			fileID = update.Message.Voice.FileID
		} else if update.Message.Audio != nil {
			fileID = update.Message.Audio.FileID
		} else if update.Message.Document != nil {
			// If user sends a file (e.g., .ogg/.mp3/.wav), try it as well
			fileID = update.Message.Document.FileID
		} else {
			continue
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.api.Send(typing)

		tgFile, err := bot.api.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			bot.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Failed to get file: %v", err)))
			continue
		}
		tmpDir := os.TempDir()
		base := filepath.Base(tgFile.FilePath)
		if base == "." || base == "/" || base == "" {
			base = fmt.Sprintf("tg_%d.oga", time.Now().UnixNano())
		}
		srcPath := filepath.Join(tmpDir, base)
		wavPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ".wav"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Downloading")
		bot.api.Send(msg)
		fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.api.Token, tgFile.FilePath)
		if err := downloadFile(fileURL, srcPath); err != nil {
			bot.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Download error: %v", err)))
			continue
		}
		defer os.Remove(srcPath)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Converting")
		bot.api.Send(msg)
		// Convert to WAV
		if err := convertToWav(srcPath, wavPath); err != nil {
			bot.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Convert error: %v", err)))
			continue
		}
		defer func() {
			os.Remove(wavPath)
			os.Remove(wavPath + ".txt")
		}()

		// Detect language from Telegram (optional; Whisper can auto)
		lang := ""
		if update.Message.From != nil && update.Message.From.LanguageCode != "" {
			// Telegram sends like "en", "ru", "uk", etc. Whisper expects ISO639-1—this is fine.
			lang = update.Message.From.LanguageCode
		}

		// Run Whisper with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		whisperBin := "./whisper/whisper-cli"
		//modelPath := "./whisper/ggml-medium.bin"
		modelPath := "./whisper/ggml-small.bin"

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Recognition")
		bot.api.Send(msg)

		_ = lang
		text, err := runWhisper(ctx, whisperBin, modelPath, wavPath, "", false)
		if err != nil {
			bot.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Whisper error: %v", err)))
			continue
		}

		if text == "" {
			text = "(no speech detected)"
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		reply.ReplyToMessageID = update.Message.MessageID
		bot.api.Send(reply)

		command := strings.ToLower(text)
		switch command {
		case "turn on the lamp.":
			msg.Text = bot.SendReq("VERTEX", "LAMP:ON")
		case "turn off the lamp.":
			msg.Text = bot.SendReq("VERTEX", "LAMP:OFF")
		default:
			continue
		}

		bot.api.Send(msg)

	}
}

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}
