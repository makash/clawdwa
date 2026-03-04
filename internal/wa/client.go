package wa

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	_ "modernc.org/sqlite"

	"github.com/makash/clawdwa/internal/config"
)

// Client wraps a whatsmeow client.
type Client struct {
	wa *whatsmeow.Client
}

// New opens the whatsmeow SQLite store and returns a Client.
func New() (*Client, error) {
	path := config.WADBPath()

	if err := os.MkdirAll(config.Dir(), 0700); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	container := sqlstore.NewWithDB(db, "sqlite3", waLog.Noop)
	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("upgrade wa db: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}

	waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)

	return &Client{wa: waClient}, nil
}

// Connect connects to WhatsApp. If not authenticated, shows a QR code.
// Blocks until authenticated and connected.
func (c *Client) Connect(ctx context.Context) error {
	if c.wa.Store.ID == nil {
		// Not logged in — show QR code.
		qrChan, err := c.wa.GetQRChannel(ctx)
		if err != nil {
			return fmt.Errorf("get qr channel: %w", err)
		}

		if err := c.wa.Connect(); err != nil {
			return fmt.Errorf("connect: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("QR login event:", evt.Event)
				if evt.Event == "success" {
					break
				}
			}
		}
	} else {
		if err := c.wa.Connect(); err != nil {
			return fmt.Errorf("connect: %w", err)
		}
	}

	return nil
}

// AddMessageHandler registers a handler for incoming text messages.
// It also registers a reconnect handler for disconnection events.
func (c *Client) AddMessageHandler(fn func(sender, chatJID, msgID, text string)) {
	c.wa.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			var text string
			if v.Message == nil {
				return
			}
			text = v.Message.GetConversation()
			if text == "" {
				if ext := v.Message.GetExtendedTextMessage(); ext != nil {
					text = ext.GetText()
				}
			}
			if text == "" {
				// Not a text message (image, video, etc.) — skip.
				return
			}
			fn(v.Info.Sender.String(), v.Info.Chat.String(), v.Info.ID, text)

		case *events.Disconnected:
			// Attempt reconnect on disconnect.
			_ = c.wa.Connect()
		}
	})
}

// SendText parses jid and sends a text message.
func (c *Client) SendText(ctx context.Context, jid, text string) error {
	parsedJID, err := types.ParseJID(jid)
	if err != nil {
		return fmt.Errorf("parse jid %q: %w", jid, err)
	}

	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	_, err = c.wa.SendMessage(ctx, parsedJID, msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

// Disconnect disconnects the WhatsApp client.
func (c *Client) Disconnect() {
	c.wa.Disconnect()
}

// Group holds basic info about a WhatsApp group.
type Group struct {
	JID  string
	Name string
}

// GetGroups returns all joined WhatsApp groups.
func (c *Client) GetGroups(ctx context.Context) ([]Group, error) {
	groups, err := c.wa.GetJoinedGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("get joined groups: %w", err)
	}
	result := make([]Group, 0, len(groups))
	for _, g := range groups {
		result = append(result, Group{
			JID:  g.JID.String(),
			Name: g.Name,
		})
	}
	return result, nil
}

// SendOnce creates a client, connects, sends text to cfg.GroupJID, then disconnects.
// Used by `clawdwa send`.
func SendOnce(cfg *config.Config, text string) error {
	c, err := New()
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}
	defer c.Disconnect()

	ctx := context.Background()
	if err := c.Connect(ctx); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return c.SendText(ctx, cfg.GroupJID, text)
}
