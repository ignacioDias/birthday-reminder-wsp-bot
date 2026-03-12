# Birthday Reminder

A WhatsApp-based birthday reminder application written in Go. The app checks daily for birthdays and sends automated WhatsApp messages to notify you about upcoming celebrations.

## Features

- 📅 Add, list, and remove birthdays via CLI
- 📱 Sends WhatsApp messages automatically on birthdays
- 🔐 Secure WhatsApp authentication using QR code
- 💾 Persistent storage of birthdays in JSON format
- 🚀 Lightweight and easy to set up as a startup program

## Prerequisites

- Go 1.25.0 or higher
- WhatsApp account
- Linux/Unix system (for startup automation)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd birthdayreminder
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build ./cmd/birthdayreminder
```

## Setup

### Initial Configuration

1. **First-time WhatsApp login:**
   ```bash
   ./birthdayreminder
   ```
   - A QR code will appear in your terminal
   - Scan it with your WhatsApp mobile app (Settings → Linked Devices → Link a Device)
   - Wait for successful authentication
   - The session will be saved for future use

2. **Set your WhatsApp number** (to receive birthday notifications):
   ```bash
   ./birthdayreminder --num 5493515234567
   ```
   Replace with your WhatsApp number in international format (country code + number, no + or spaces)

### Adding as a Startup Program

To check birthdays automatically every day, add the application to your system's startup programs

## Usage

### Add a Birthday
```bash
./birthdayreminder --add Juan 01/05
```
Format: `--add <name> <DD/MM>`

### List All Birthdays
```bash
./birthdayreminder --list
```

### Remove a Birthday
```bash
./birthdayreminder --remove Juan
```
⚠️ Note: This removes all birthdays for the given name. Names are case-sensitive.

### Set/Update Phone Number
```bash
./birthdayreminder --num 5493515234567
```

### Check Birthdays (Manual)
```bash
./birthdayreminder
```
Runs without flags to check today's birthdays and send messages if any match.

### Display Help
```bash
./birthdayreminder --help
```

## Data Storage

The application stores data in `~/Documents/birthdayapp/`:
- `birthdays.json` - List of all birthdays
- `number.json` - WhatsApp number for notifications
- `session.db` - WhatsApp session data (SQLite)

## How It Works

1. On startup, the app loads stored birthdays and your configured phone number
2. It checks if today's date matches any stored birthdays
3. If matches are found, it connects to WhatsApp and sends a message: "Hoy es cumpleaños de [Name]"
4. The WhatsApp session persists, so you only need to scan the QR code once

## Dependencies

- [whatsmeow](https://github.com/tulir/whatsmeow) - WhatsApp Web API
- [qrterminal](https://github.com/mdp/qrterminal) - QR code generation in terminal
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for session storage

## Troubleshooting

### QR Code Not Appearing
- Ensure your terminal supports UTF-8 encoding
- Try clearing the session: `rm -rf ~/Documents/birthdayapp/session.db`
- Run the application again

### Messages Not Sending
- Verify your WhatsApp number is set correctly with `--num`
- Ensure the WhatsApp session is still valid (you may need to re-authenticate)
- Check that birthdays are added in the correct format (DD/MM)

### Permission Issues
- Ensure the application has write permissions to `~/Documents/birthdayapp/`
