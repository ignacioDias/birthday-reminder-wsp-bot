package cli

import (
	"birthdayreminder/internal/wsp"
	"flag"
	"time"

	"go.mau.fi/whatsmeow"
)

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

func NewCLI(args []string, number string) *CLI {
	return &CLI{
		args:      args,
		client:    nil,
		birthdays: loadBirthdays(),
		number:    number,
	}
}

func (cli *CLI) Run() error {
	if len(cli.args) == 0 {
		return cli.checkBirthdays()
	}
	add := flag.String("add", "", "Add a birthday: Name, DD/MM. Example: --add Juan, 01/05")
	list := flag.Bool("list", false, "List of all birthdays: --list")
	num := flag.String("num", "", "Set number to send message: --num 5493515234567")
	remove := flag.String("remove", "", "Removes all birthdays from a given name. Careful with capital letters. Example: --remove Juan")
	help := flag.Bool("help", false, "List of commands: --help")

	flag.Parse()

	switch {
	case *add != "":

	case *list:
	case *remove != "":
	case *num != "":
		cli.number = *num
	case *help:
		printHelp()
	}
	return nil
}

func printHelp() {

}

func loadBirthdays() []Birthday {
	return nil
}

func (cli *CLI) saveBirthdays() error {
	return nil
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
