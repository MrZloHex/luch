package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	log "log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"luch/pkg/audioconv"
	"luch/pkg/stt"
)

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

func (bot *Bot) GetTextOrVoice(update tgbotapi.Update) (string, error) {
	var fileID string
	if update.Message.Voice != nil {
		fileID = update.Message.Voice.FileID
	} else if update.Message.Audio != nil {
		fileID = update.Message.Audio.FileID
	} else if update.Message.Document != nil {
		fileID = update.Message.Document.FileID
	} else {
		return update.Message.Text, nil
	}

	typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
	bot.SendBot(typing)

	tgFile, err := bot.GetFileBot(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		bot.SendBot(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Failed to get file: %v", err)))
		return "", err
	}
	tmpDir := os.TempDir()
	base := filepath.Base(tgFile.FilePath)
	if base == "." || base == "/" || base == "" {
		base = fmt.Sprintf("tg_%d.ogg", time.Now().UnixNano())
	}
	srcPath := filepath.Join(tmpDir, base)
	wavPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ".wav"

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.GetToken(), tgFile.FilePath)
	if err := downloadFile(fileURL, srcPath); err != nil {
		bot.SendBot(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Download error: %v", err)))
		return "", err
	}
	defer os.Remove(srcPath)

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


	res, err := bot.tr.TranscribePCM(context.Background(), pcm, stt.Options{
		Language:        "auto",
		TranslateToEn:   false,
		Threads:         0, // auto (NumCPU)
		TokenTimestamps: false,
		MaxSegmentChars: 0,
		MaxTokens:       0,
		BeamSize:        0, // >0 enables beam search
		// Optional slicing:
		Offset:   0 * time.Second, // start at 0
		Duration: 0,               // full length
	})
	if err != nil {
		log.Error("transcribe: %v", err)
	}

	log.Info("STT: ", "res", res.Text)

	return res.Text, nil
}
