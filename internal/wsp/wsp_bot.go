package wsp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func NewBot() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "DEBUG", false)
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dbPath := filepath.Join(wd, "internal", "database", "session.db")
	container, err := sqlstore.New(ctx, "sqlite3", "file:"+dbPath+"?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	return client
}

func eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func SendMessage(client *whatsmeow.Client, phone string, message string) error {
	jid, err := types.ParseJID(phone + "@s.whatsapp.net")
	if err != nil {
		return err
	}

	_, err = client.SendMessage(context.Background(), jid, &waProto.Message{
		Conversation: proto.String(message),
	})
	return err
}
