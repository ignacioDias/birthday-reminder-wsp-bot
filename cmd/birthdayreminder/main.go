package main

import (
	"birthdayreminder/internal/wsp"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bot := wsp.NewBot()
	for !bot.IsConnected() {
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Enviando mensaje...")
	if err := wsp.SendMessage(bot, "num", "Hola, soy un humano!"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Mensaje enviado!")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	bot.Disconnect()
}
