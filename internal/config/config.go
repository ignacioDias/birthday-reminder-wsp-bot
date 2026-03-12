package config

import (
	"birthdayreminder/internal/models"
	"bufio"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homeDir, "Documents", "birthdayapp", "number.txt")

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", scanner.Err()
}

func readFromBirthdaysFile() ([]models.Birthday, error) {
	return readFromFile[[]models.Birthday]("birthdays")
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

func fileExists(path string) bool {
	_, error := os.Stat(path)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
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

	dir := filepath.Join(homeDir, "Documents", "birthdayapp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file := filepath.Join(dir, fileName+".json")
	return os.WriteFile(file, data, 0644)
}

func SaveNumber(number string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(homeDir, "Documents", "birthdayapp", "number.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(number + "\n")
	return err
}
