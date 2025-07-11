package card

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// 读取保存抽卡记录
type CardStore interface {
	Save(CardSaveType) error
	Load() (CardSaveType, error)
}

// 保存抽卡记录到系统文件中
type CardStoreToFile struct {
	FilePath string
}

// 保存数据到文件,不存在时创建
func (c *CardStoreToFile) Save(data CardSaveType) error {
	if c.FilePath == "" {
		return errors.New("file path is empty")
	}
	if !data.Valid() {
		return errors.New("data is invalid")
	}
	if err := os.MkdirAll(filepath.Dir(c.FilePath), os.ModePerm); err != nil {
		return err
	}
	file, err := os.OpenFile(c.FilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data.String())
	if err != nil {
		return err
	}
	return nil
}
func (c *CardStoreToFile) Load() (CardSaveType, error) {
	if c.FilePath == "" {
		return CardSaveType{}, errors.New("file path is empty")
	}
	file, err := os.Open(c.FilePath)
	if err != nil {
		return CardSaveType{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return CardSaveType{}, err
	}
	card := CardSaveType{}
	err = card.Parse(string(data))
	return card, err
}
