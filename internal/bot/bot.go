package bot

import (
	"fmt"
	stdlog "log"
	log "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"luch/pkg/protocol"

	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"luch/pkg/audioconv"
	"luch/pkg/stt"
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

func (bot *Bot) SendReq(parts ...string) string {
	resp, err := bot.ptcl.Send(parts...)
	if err != nil {
		return fmt.Sprintf("Failed to send request: %s", err.Error())
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


func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.CallbackQuery != nil {
			bot.proccessInlineKeyboard(update)
		}

		if update.Message == nil {
			continue
		}

		log.Debug("Got smth", "from", update.Message.From.UserName, "text", update.Message.Text)

		switch {
		case update.Message.IsCommand():
			bot.processCmd(update)
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
			base = fmt.Sprintf("tg_%d.ogg", time.Now().UnixNano())
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
		
		pcm, err := audioconv.ConvertFileToPCM16k(context.Background(), srcPath, audioconv.Options{})
		if err != nil {
			log.Error("convert audio: %v", err)
		}
		if len(pcm) == 0 {
			log.Error("no audio samples after conversion")
		}
		defer func() {
			os.Remove(wavPath)
			os.Remove(wavPath + ".txt")
		}()


		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Recognition")
		bot.api.Send(msg)

		tr, err := stt.NewTranscriber("third_party/whisper.cpp/models/ggml-base.bin")
		if err != nil {
			log.Error("load model: %v", err)
		}
		defer tr.Close()

		res, err := tr.TranscribePCM(context.Background(), pcm, stt.Options{
			Language:        "auto",
			TranslateToEn:   false,
			Threads:         0,            // auto (NumCPU)
			TokenTimestamps: false,
			MaxSegmentChars: 0,
			MaxTokens:       0,
			BeamSize:        0,        // >0 enables beam search
			// Optional slicing:
			Offset:   0 * time.Second,     // start at 0
			Duration: 0,                    // full length
		})
		if err != nil {
			log.Error("transcribe: %v", err)
		}

		// 4) Print results
		fmt.Printf("Language: %s\n", res.Language)
		fmt.Println("Segments:")
		for i, s := range res.Segments {
			fmt.Printf("  %2d  [%6.2f .. %6.2f]  %s\n", i, s.StartSec, s.EndSec, s.Text)
		}
		fmt.Println("\nFULL TEXT:")
		fmt.Println(res.Text)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, res.Text)
		reply.ReplyToMessageID = update.Message.MessageID
		bot.api.Send(reply)

	}
}

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}
