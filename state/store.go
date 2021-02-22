package state

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	bot "github.com/genvmoroz/bot-engine/api"
	"github.com/genvmoroz/bot-engine/processor"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

const (
	Store = "/store"
	End   = "/end"
)

type storeStateProcessor struct {
	tgBot            bot.Client
	chatID           int64
	fileStorePath    string
	totalStoredFiles uint32
}

func NewStoreStateProcessor(tbBot bot.Client, chatID int64, fileStorePath string) processor.StateProcessor {
	return &storeStateProcessor{
		tgBot:         tbBot,
		chatID:        chatID,
		fileStorePath: fileStorePath,
	}
}

func (p *storeStateProcessor) Process(updateChan <-chan tgBotApi.Update) error {
	if updateChan == nil {
		return errors.New("updateChan cannot be nil")
	}

	p.totalStoredFiles = 0
	if err := p.tgBot.Send("You're in the store state.", p.chatID); err != nil {
		return fmt.Errorf("failed to send the message: %w", err)
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
				msg = fmt.Sprintf("%s, %d files were stored", msg, p.totalStoredFiles)
			} else {
				msg = fmt.Sprintf("%s, 1 file was stored", msg)
			}
			return p.tgBot.Send(msg, p.chatID)
		}

		if err := p.resolveFilesAndStore(*update.Message); err != nil {
			return fmt.Errorf("failed to resolve files amd store: %w", err)
		}
	}
}

func (p *storeStateProcessor) resolveFilesAndStore(message tgBotApi.Message) error {
	var storedCount uint32
	if message.Document != nil {
		if err := p.storeFile(message.Document.FileID, message.Document.FileName); err != nil {
			return fmt.Errorf("failed to store the document: %w", err)
		}
		storedCount++
	}
	if message.Photo != nil && len(*message.Photo) != 0 {
		largest := getLargestPhotoByFileSize(*message.Photo...)
		if largest == nil {
			return fmt.Errorf("largest photo cannot be nil")
		}
		if err := p.storeFile(largest.FileID, uuid.New().String()+".jpg"); err != nil {
			return fmt.Errorf("failed to store the photo: %w", err)
		}
		storedCount++
	}

	p.totalStoredFiles += storedCount

	if storedCount == 0 {
		if err := p.tgBot.Send("0 files were stored, in order to store some files you should upload some files.", p.chatID); err != nil {
			return fmt.Errorf("failed to send the message: %w", err)
		}
	} else if err := p.tgBot.Send(fmt.Sprintf("%d files stored.", storedCount), p.chatID); err != nil {
		return fmt.Errorf("failed to send the message: %w", err)
	}

	return nil
}

func (p *storeStateProcessor) storeFile(fileID, fileName string) error {
	if strings.Trim(fileID, " ") == "" {
		return errors.New("file id should be provided")
	}
	if strings.Trim(fileName, " ") == "" {
		return errors.New("file name should be provided")
	}

	if err := p.downloadFile(fileID, fileName); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	creationTimeAsString, err := p.getModTimeAsString(fileName)
	if err != nil {
		return fmt.Errorf("failed to get creation time as string: %w", err)
	}
	if err = p.saveFileInDirectory(creationTimeAsString, fileName); err != nil {
		return fmt.Errorf("failed to save file in directory: %w", err)
	}

	return nil
}

func (p *storeStateProcessor) saveFileInDirectory(directory, fileName string) error {
	if strings.Trim(fileName, " ") == "" {
		return errors.New("file name should be provided")
	}
	if strings.Trim(directory, " ") == "" {
		return errors.New("directory name should be provided")
	}

	if p.fileStorePath[len(p.fileStorePath)-1] != '/' {
		p.fileStorePath += "/"
	}

	rootPath := fmt.Sprintf("%s%s", p.fileStorePath, directory)
	if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to mkdir all directories: %w", err)
	}

	return os.Rename(fileName, fmt.Sprintf("%s/%s", rootPath, fileName))
}

func (p *storeStateProcessor) downloadFile(fileID, fileName string) error {
	if strings.Trim(fileID, " ") == "" {
		return errors.New("file id should be provided")
	}
	if strings.Trim(fileName, " ") == "" {
		return errors.New("file name should be provided")
	}

	content, err := p.tgBot.DownloadFile(fileID)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if err = os.WriteFile(fileName, content, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (p *storeStateProcessor) getModTimeAsString(fileName string) (string, error) {
	if strings.Trim(fileName, " ") == "" {
		return "", errors.New("file name should be provided")
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read stats of file %s: %w", fileName, err)
	}

	return fileInfo.ModTime().Format("2006-01-02"), nil
}

func getLargestPhotoByFileSize(photos ...tgBotApi.PhotoSize) *tgBotApi.PhotoSize {
	if len(photos) == 0 {
		return nil
	}

	largest := &photos[0]
	for i := 1; i < len(photos); i++ {
		if photos[i].FileSize > largest.FileSize {
			largest = &photos[i]
		}
	}

	return largest
}
