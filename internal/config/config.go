package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadBirthdays() (any, error) {
	birthdays, err := readFromFile("birthdays")
	if err != nil {
		return nil, err
	}
	return birthdays, nil
}
func LoadNumber() (string, error) {
	number, err := readFromFile("number")
	if err != nil {
		return "", err
	}
	return number, nil
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
