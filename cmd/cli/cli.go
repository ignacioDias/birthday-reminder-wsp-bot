package cli

import (
	"birthdayreminder/internal/config"
	"birthdayreminder/internal/models"
	"birthdayreminder/internal/wsp"
	"errors"
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
)

type CLI struct {
	args      []string
	client    *whatsmeow.Client
	birthdays []models.Birthday
	number    string
}

func NewCLI(args []string) (*CLI, error) {
	number, err := config.LoadNumber()
	if err != nil {
		return nil, err
	}
	birthdays, err := config.LoadBirthdays()
	if err != nil {
		return nil, err
	}
	return &CLI{
		args:      args,
		client:    nil,
		birthdays: birthdays,
		number:    number,
	}, nil
}

func (cli *CLI) Run() error {
	if len(cli.args) == 0 {
		return cli.checkBirthdays()
	}
	add := flag.Bool("add", false, "Add a birthday: Name, DD/MM. Example: --add Juan 01/05")
	list := flag.Bool("list", false, "List of all birthdays: --list")
	num := flag.Bool("num", false, "Set number to send message. Needed for start: --num 5493515234567")
	remove := flag.Bool("remove", false, "Removes all birthdays from a given name. Careful with capital letters. Example: --remove Juan")
	help := flag.Bool("help", false, "List of commands: --help")

	flag.Parse()
	args := flag.Args()
	switch {
	case *add:
		if len(args) < 2 {
			return errors.New("Error: requires <name> and <DD/MM>")
		}
		if err := cli.addBirthday(args[0], args[1]); err != nil {
			return err
		}
		return config.SaveFile(cli.birthdays, "birthdays")
	case *list:
		cli.listBirthdays()
	case *remove:
		if len(args) < 1 {
			return errors.New("Error: requires <name>")
		}
		cli.removeBirthdays(args[0])
		return config.SaveFile(cli.birthdays, "birthdays")
	case *num:
		if len(args) < 1 {
			return errors.New("Error: requires <number>")
		}
		cli.number = args[0]
		config.SaveFile(cli.number, "number")
	case *help:
		printHelp()
	}
	return nil
}

func (cli *CLI) addBirthday(name string, date string) error {
	parts := strings.Split(date, "/")
	if len(parts) != 2 {
		return errors.New("Invalid date")
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return errors.New("Invalid day")
	}

	monthInt, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("Invalid month")
	}
	month := time.Month(monthInt)
	birthday := models.Birthday{
		Name:  name,
		Day:   day,
		Month: month,
	}

	cli.birthdays = append(cli.birthdays, birthday)
	return nil
}

func printHelp() {
	fmt.Println("Birthday Reminder CLI")
	fmt.Println("\nUsage:")
	fmt.Println("\nOn first run, execute without commands to log in. Then you can add it as a startup program.")
	fmt.Println("\n Then, you add the number to send the reminders (I suggest to use your own number).")
	fmt.Println("  birthdayreminder [flags] [arguments]")
	fmt.Println("\nFlags:")
	fmt.Println("  --add <name> <DD/MM>    Add a birthday")
	fmt.Println("                          Example: --add Juan 01/05")
	fmt.Println("  --list                  List all birthdays")
	fmt.Println("  --remove <name>         Remove all birthdays for a given name")
	fmt.Println("                          Example: --remove Juan")
	fmt.Println("  --num <number>          Set WhatsApp number to send messages")
	fmt.Println("                          Example: --num 5493515234567")
	fmt.Println("  --help                  Display this help message")
	fmt.Println("\nDefault behavior (no flags):")
	fmt.Println("  Check today's birthdays and send WhatsApp messages if any match")
}

func (cli *CLI) listBirthdays() {
	for _, birthday := range cli.birthdays {
		fmt.Printf("%s: %d/%s\n", birthday.Name, birthday.Day, birthday.Month.String())
	}
}

func (cli *CLI) removeBirthdays(name string) {
	for i := len(cli.birthdays) - 1; i >= 0; i-- {
		if cli.birthdays[i].Name == name {
			cli.birthdays = slices.Delete(cli.birthdays, i, i+1)
		}
	}
}

func (cli *CLI) checkBirthdays() error {
	var err error
	cli.client, err = wsp.NewBot()
	if err != nil {
		return err
	}
	for _, elem := range cli.birthdays {
		if elem.Month == time.Now().Month() && elem.Day == time.Now().Day() {
			if err := wsp.SendMessage(cli.client, cli.number, generateMessage(elem.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateMessage(name string) string {
	return "Hoy es cumpleaños de " + name
}

func (cli *CLI) KillBot() {
	if cli.client != nil {
		cli.client.Disconnect()
	}
}
