package main

import (
	"fmt"
	"os"
)

// ============================================
// 1. INTERFACE — just the contract
// ============================================

type Notification interface {
	Send(message string) string
}

// ============================================
// 2. STRUCTS — real implementation, hidden complexity
// ============================================

type EmailNotification struct {
	SMTPServer string
	Port       int
}
type WhatsappNotification struct {
	AccessToken string
	PhoneID     string
}

// ============================================
// 3. METHOD IMPLEMENTATIONS — each struct does its own thing
// ============================================

func (e EmailNotification) Send(message string) string {
	// implement email send

	return fmt.Sprintf("email sent via", e.SMTPServer, e.Port)
}

func (w WhatsappNotification) Send(message string) string {
	// implement the logic to send whatsapp message here

	return fmt.Sprintf("whatsapp message sent to phone id", w.PhoneID)
}

// ============================================
//  4. FACTORY — this is just a regular function
//     it knows all struct names
//     it knows how to init each one
//     it knows where to get credentials
//     caller knows NONE of this
//
// ============================================
func NewNotificationFactory(channel string) Notification {
	switch channel {
	case "email":
		return EmailNotification{
			SMTPServer: "smtp.gmail.com",
			Port:       587,
		}
	case "whatsapp":
		return WhatsappNotification{
			AccessToken: os.Getenv("ACCESS_TOKEN"),
			PhoneID:     "my_phone_id",
		}
	default:
		panic("unknown channel: " + channel)
	}
}

// ============================================
// 5. CALLER — knows nothing except channel name and Send()
// ============================================

func main() {
	channels := []string{"email", "whatsapp"}

	for _, channel := range channels {
		n := NewNotificationFactory(channel)
		res := n.Send("your order has shipped")
		fmt.Println(res)
	}
}
