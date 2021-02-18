package processor

import (
	"errors"
	"fmt"
	"log"
	"os"

	"ff-bot/bot"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type storeStateProcessor struct {
	baseStateProcessor
	fileStorePath    string
	totalStoredFiles uint32
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
				msg = fmt.Sprintf("%s, %d files stored", msg, p.totalStoredFiles)
			} else {
				msg = fmt.Sprintf("%s, 1 file stored", msg)
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
		if err := p.storeDocument(*message.Document); err != nil {
			return fmt.Errorf("failed to store document: %w", err)
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

func (p *storeStateProcessor) storeDocument(document tgBotApi.Document) error {
	if err := p.downloadDocument(document); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	creationTimeAsString, err := p.getModTimeAsString(document.FileName)
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

func (p *storeStateProcessor) getModTimeAsString(fileName string) (string, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read stats of file %s: %w", fileName, err)
	}

	return fileInfo.ModTime().Format("2006-01-02"), nil
}
