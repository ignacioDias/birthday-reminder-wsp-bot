package cli

import (
	"birthdayreminder/internal/wsp"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"time"

	"go.mau.fi/whatsmeow"
)

var FILE string = "./save.json"

type Birthday struct {
	date time.Time
	name string
}

type CLI struct {
	args      []string
	client    *whatsmeow.Client
	birthdays []Birthday
	number    string
}

func NewCLI(args []string) *CLI {
	return &CLI{
		args:      args,
		client:    nil,
		birthdays: loadBirthdays(),
		number:    loadNumber(),
	}
}

func (cli *CLI) Run() error {
	if len(cli.args) == 0 {
		return cli.checkBirthdays()
	}
	add := flag.String("add", "", "Add a birthday: Name, DD/MM. Example: --add Juan, 01/05")
	list := flag.Bool("list", false, "List of all birthdays: --list")
	num := flag.String("num", "", "Set number to send message. Needed for start: --num 5493515234567")
	remove := flag.String("remove", "", "Removes all birthdays from a given name. Careful with capital letters. Example: --remove Juan")
	help := flag.Bool("help", false, "List of commands: --help")

	flag.Parse()

	switch {
	case *add != "":
		return cli.saveBirthdays()
	case *list:
		cli.listBirthdays()
	case *remove != "":
		cli.removeBirthdays(*remove)
		return cli.saveBirthdays()
	case *num != "":
		cli.number = *num
	case *help:
		printHelp()
	}
	return nil
}

func printHelp() {

}

func (cli *CLI) listBirthdays() {
	for _, birthday := range cli.birthdays {
		fmt.Printf("%s: %d/%d\n", birthday.name, birthday.date.Day(), birthday.date.Month())
	}
}

func (cli *CLI) removeBirthdays(name string) {
	for i := len(cli.birthdays) - 1; i >= 0; i-- {
		if cli.birthdays[i].name == name {
			cli.birthdays = slices.Delete(cli.birthdays, i, i+1)
		}
	}
}

func loadBirthdays() ([]Birthday, error) {
	birthdays, err := readFromFile()
	if err != nil {
		return nil, err
	}
	return birthdays, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func (cli *CLI) saveBirthdays() error {
	data, err := json.MarshalIndent(cli.birthdays, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(FILE, data, 0644)
	return err
}

func (cli *CLI) checkBirthdays() error {
	cli.client = wsp.NewBot()
	for _, elem := range cli.birthdays {
		if elem.date.Month() == time.Now().Month() && elem.date.Day() == time.Now().Day() {
			if err := wsp.SendMessage(cli.client, cli.number, generateMessage(elem.name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateMessage(name string) string {
	return "Hoy es cumpleaños de " + name
}
