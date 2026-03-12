package wsp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func NewBot() (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", false)
	ctx := context.Background()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(homeDir, "Documents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, "birthdayapp", "session.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	container, err := sqlstore.New(ctx, "sqlite3", "file:"+dbPath+"?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "DEBUG", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		connected := make(chan struct{})
		client.AddEventHandler(func(evt interface{}) {
			switch evt.(type) {
			case *events.Connected:
				close(connected)
			}
		})

		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		for evt := range qrChan {
			switch evt.Event {
			case "code":
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			case "success":
				<-connected
				fmt.Println("Logged in successfully! Run again to start the bot.")
				time.Sleep(2 * time.Second)
				client.Disconnect()
				os.Exit(0)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		for !client.IsConnected() {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return client, nil
}

func SendMessage(client *whatsmeow.Client, phone string, message string) error {
	jid, err := types.ParseJID(phone + "@s.whatsapp.net")
	if err != nil {
		return err
	}

	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: proto.String(message),
	})
	return err
}
