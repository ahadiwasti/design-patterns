# Factory Method Pattern

## What is it?

A function that takes a decision input and returns an interface, hiding the concrete struct underneath.

The caller never knows which struct it got — it just knows what it can do.

```go
func NewNotification(channel string) Notification
```

This one line is the factory. That is it.

---

## The Three Moving Parts

### 1. Interface — the contract
Defines what the caller can do. Only method signatures, no fields, no logic.

```go
type Notification interface {
    Send(message string) string
}
```

### 2. Structs — the real implementation
Each struct implements the interface in its own way. Different protocols, different APIs, different credentials — all hidden inside.

```go
type EmailNotification struct {
    SMTPServer string
    Port       int
}

func (e EmailNotification) Send(message string) string {
    // SMTP logic, TLS, headers — all hidden here
    return "[EMAIL] " + message
}

type WhatsAppNotification struct {
    AccessToken string
    PhoneID     string
}

func (w WhatsAppNotification) Send(message string) string {
    // HTTP POST to Meta API — completely different, caller doesnt care
    return "[WHATSAPP] " + message
}
```

### 3. Factory — the decision maker
Knows all the struct names, knows how to initialize them, knows where to get credentials. The caller knows none of this.

```go
func NewNotification(channel string) Notification {
    switch channel {
    case "email":
        return EmailNotification{
            SMTPServer: os.Getenv("SMTP_SERVER"),
            Port:       587,
        }
    case "whatsapp":
        return WhatsAppNotification{
            AccessToken: os.Getenv("WA_TOKEN"),
            PhoneID:     os.Getenv("WA_PHONE_ID"),
        }
    default:
        panic("unknown channel: " + channel)
    }
}
```

---

## Why the return type must be the Interface

Without interface return type — concrete type leaks to caller:

```go
// bad — returns concrete type
func NewNotification(channel string) EmailNotification

n := NewNotification("email")
n.SMTPServer = "something"  // caller knows too much, tightly coupled
```

With interface return type — caller is sealed off from implementation:

```go
// good — returns interface
func NewNotification(channel string) Notification

n := NewNotification("email")  // n is Notification, caller has no idea what is underneath
n.Send("hello")                // this is all caller can do
```

The interface return type is the seal. It enforces "I will give you something that can Send — what exactly it is, none of your business."

---

## What Each Part Is Responsible For

| Part | Role | When |
|---|---|---|
| Interface | contract / promise | compile time |
| Factory | picks and builds the right struct | runtime |
| Struct | actually executes the logic | runtime |

The interface is invisible at runtime. It did its job at compile time — it made sure every struct that claims to be a `Notification` actually has a `Send` method with the correct signature.

---

## Why Not Just One Struct?

A single struct with if/else works but does not scale:

```go
// this grows forever
func (e Notification) Send(message string, channel string) string {
    if channel == "email"    { ... }
    if channel == "whatsapp" { ... }
    if channel == "telegram" { ... }  // touch existing code every time
}
```

Every new channel means editing the same function. You risk breaking email while adding Telegram.

With the factory pattern, adding a new channel means:
- Add a new struct
- Add one case to the factory
- Everything else stays untouched

---

## The Full Flow

```
caller: NewNotification("whatsapp")
              ↓
        factory hits case "whatsapp"
              ↓
        builds WhatsAppNotification{AccessToken: "...", PhoneID: "..."}
              ↓
        returns it as Notification interface type
              ↓
        n.Send("hello")
              ↓
        Go sees actual type is WhatsAppNotification
              ↓
        calls WhatsAppNotification.Send() — not Email, not SMS
```

---

## Scaling Further — Registry Pattern

When the factory switch itself becomes a problem, replace it with a registry map:

```go
var registry = map[string]func() Notification{
    "email":    func() Notification { return EmailNotification{SMTPServer: "smtp.gmail.com"} },
    "whatsapp": func() Notification { return WhatsAppNotification{AccessToken: "xyz"} },
}

// factory never changes again
func NewNotification(channel string) Notification {
    constructor, exists := registry[channel]
    if !exists {
        panic("unknown channel: " + channel)
    }
    return constructor()
}

// adding new channel — just add to registry, factory untouched
func init() {
    registry["telegram"] = func() Notification {
        return TelegramNotification{BotToken: "abc"}
    }
}
```

---

## Real Production Fallback Use Cases

The factory becomes even more powerful when combined with health checks.
The caller never knows if it is talking to the primary or fallback — it just works.

| Use Case | Primary | Fallback |
|---|---|---|
| Payments | Stripe | PayPal |
| SMS | Twilio | AWS SNS |
| AI Models | GPT-4 | Claude |
| Storage | S3 | GCS |
| Email | SendGrid | Mailgun |

The factory checks health at decision time and returns whichever is alive.
Your business logic never changes — resilience is handled at the factory level.

---

## Key Principles This Pattern Follows

**Open/Closed Principle** — open for extension, closed for modification. Add new channels by adding new structs, never by editing existing code.

**Information Hiding** — the caller has one concern: `Send`. Everything else (HTTP calls, SMTP, API tokens, retries) is hidden inside the struct.

**Single Responsibility** — each struct owns exactly one channel's logic. The factory owns exactly one job: deciding which struct to build.

---

## Go vs Python

| | Go | Python |
|---|---|---|
| Abstract type | `interface` | `ABC` with `@abstractmethod` |
| Implementation check | compile time | runtime |
| Keyword to implement | none — implicit | none — implicit |
| Factory style | standalone function | function or class method |
| Dispatch | `switch` statement | `dict` lookup |

In Go, if your struct does not implement the interface the code will not compile. In Python it fails at runtime. Same pattern, different enforcement timing.