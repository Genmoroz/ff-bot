package processor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"ff-bot/bot"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type storeStateProcessor struct {
	baseStateProcessor
	fileStorePath    string
	totalStoredFiles int
}

func NewStoreStateProcessor(tbBot bot.Client, chatID int64, fileStorePath string) StateProcessor {
	return &storeStateProcessor{
		baseStateProcessor: newBaseStateProcessor(tbBot, chatID),
		fileStorePath:      fileStorePath,
	}
}

func (p *storeStateProcessor) Process(updateChan tgBotApi.UpdatesChannel) error {
	if updateChan == nil {
		return errors.New("updateChan cannot be nil")
	}

	p.totalStoredFiles = 0
	if err := p.tgBot.Send("You're in the store state.", p.chatID); err != nil {
		return fmt.Errorf("failed to send the message[chatID:%d]: %w", p.chatID, err)
	}
	for {
		update, ok := <-updateChan
		if !ok {
			log.Printf("updateChan is closed")
			return nil
		}

		text := update.Message.Text
		if text == End {
			msg := "End of the store state"
			if p.totalStoredFiles != 1 {
				msg = fmt.Sprintf("%s, %d files stored", msg, p.totalStoredFiles)
			} else {
				msg = fmt.Sprintf("%s, 1 file stored", msg)
			}
			return p.tgBot.Send(msg, p.chatID)
		}

		if err := p.resolveFilesAndStore(*update.Message); err != nil {

		}
	}
}

func (p *storeStateProcessor) resolveFilesAndStore(message tgBotApi.Message) error {
	var storedCount int32
	if message.Document != nil {
		if err := p.storeDocument(*message.Document); err != nil {
			return fmt.Errorf("failed to store document: %w", err)
		}
		storedCount++
	}

	if storedCount == 0 {
		if err := p.tgBot.Send("0 files were stored, in order to store some files you should upload some files.", p.chatID); err != nil {
			return fmt.Errorf("failed to send the message[chatID:%d]: %w", p.chatID, err)
		}
	}

	return nil
}

func (p *storeStateProcessor) storeDocument(document tgBotApi.Document) error {
	if err := p.downloadDocument(document); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	creationTimeAsString, err := p.getCreationTimeAsString(document.FileName)
	if err != nil {
		return fmt.Errorf("failed to get creation time as string: %w", err)
	}
	if err = p.storeFile(creationTimeAsString, document.FileName); err != nil {
		return fmt.Errorf("failed to store file: %w", err)
	}

	return nil
}

func (p *storeStateProcessor) storeFile(directory, fileName string) error {
	if p.fileStorePath[len(p.fileStorePath)-1] != '/' {
		p.fileStorePath += "/"
	}

	rootPath := fmt.Sprintf("%s%s", p.fileStorePath, directory)

	if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to mkdir all directories: %w", err)
	}

	return os.Rename(fileName, fmt.Sprintf("%s/%s", rootPath, fileName))
}

func (p *storeStateProcessor) downloadDocument(document tgBotApi.Document) error {
	content, err := p.tgBot.DownloadFile(document.FileID)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if err = os.WriteFile(document.FileName, content, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (p *storeStateProcessor) getCreationTimeAsString(fileName string) (string, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read stats of file %s: %w", fileName, err)
	}

	creationTime := time.Now()
	switch runtime.GOOS {
	case "windows":
		win32File := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		creationTime = fileTimeToTime(win32File.CreationTime.LowDateTime, win32File.CreationTime.HighDateTime)
	case "darwin":
		return "", fmt.Errorf("darwin: %+v", fileInfo.Sys())
	case "linux":
		return "", fmt.Errorf("linux: %+v", fileInfo.Sys())
	default:
		return "", fmt.Errorf("unknown GOOS: %s", runtime.GOOS)
	}

	return creationTime.Format("2006-01-02"), nil
}

func fileTimeToTime(low, high uint32) time.Time {
	var fileTime int64
	fileTime = int64(high)
	fileTime <<= 32
	fileTime += int64(low)

	return time.Unix(fileTime/10000000-11644473600, 0)
}
