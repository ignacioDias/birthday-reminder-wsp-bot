package config

import (
	"birthdayreminder/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadBirthdays() ([]models.Birthday, error) {
	birthdays, err := readFromBirthdaysFile()
	if err != nil {
		return nil, err
	}
	return birthdays, nil
}
func LoadNumber() (string, error) {
	number, err := readFromNumberFile()
	if err != nil {
		return "", err
	}
	return number, nil
}

func readFromBirthdaysFile() ([]models.Birthday, error) {
	return readFromFile[[]models.Birthday]("birthdays")
}

func readFromNumberFile() (string, error) {
	return readFromFile[string]("number")
}

func readFromFile[T any](fileName string) (T, error) {
	var result T
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return result, err
	}

	file := filepath.Join(homeDir, "Documents", "birthdayapp", fileName+".json")
	if !fileExists(file) {
		return result, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func SaveFile(dataToSave any, fileName string) error {
	data, err := json.MarshalIndent(dataToSave, "", "  ")
	if err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(homeDir, "Documents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file := filepath.Join(dir, fileName+".json")
	return os.WriteFile(file, data, 0644)
}
