# Telegram Bot Killer (Go)
This is a simple Telegram bot killer script written in Go that sends periodic messages to a specified Telegram chat. It is designed to run automatically via cron on a server, and logs all activity for monitoring purposes.

---

## 🧾 Features
- Sends messages to a Telegram bot chat using the Telegram API
- Customizable via `.env` file
- Daily log rotation stored in a `logs/` directory
- Automatically builds the binary from `kill.go` if it does not exist
- Automatically removes logs older than 30 days

---

## 📁 Project Structure
```
telegram-bot-killer/
├── kill-bot.go # Go source file
├── kill-bot # Compiled binary (auto-generated)
├── .env # Environment variable configuration
├── run.sh # Main execution script (called by cron)
└── logs-cron/
└── log-cron-YYYY-MM-DD.txt # Daily logs
```
---

## ⚙️ Configuration
### 1. Create a `.env` file in the project root directory

Example:

```env
BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
CHAT_ID=123456789
REQUEST_EVERY_TRIES=100
SLEEP_TIME=5
```

⚠️ Replace BOT_TOKEN and CHAT_ID with your actual Telegram bot token and target chat ID.

### 2. Replace ```PROJECT_DIR``` variable in run.sh with the actual full path to your project

--- 

## 🚀 Running the Bot
### 1. Run Manually (for testing)
```
bash run.sh
```
Output will be logged to:
```
logs-cron/log-cron-YYYY-MM-DD.txt
```
---

## 🧹 Log Cleanup
```
find logs-cron/ -name "log-cron-*.txt" -type f -mtime +30 -delete
```
---

## 📅 Scheduling with Cron (Linux)
To run the bot automatically at regular intervals, you can use the Linux cron scheduler.

### Step 1: Open crontab
Run the following command to edit the current user's crontab:
```
crontab -e
```

### Step 2: Add a Cron Job
Add this line at the bottom of the crontab file:
```
*/5 * * * * /bin/bash /path/to/your/project/run.sh
```
🔁 This example runs the bot every 5 minutes. Adjust the timing as needed.

Make sure to replace ```/path/to/your/project/``` with the actual full path to your project directory (```e.g. /home/username/telegram-bot/```).


## 📦 Important Notes
Make sure run.sh is executable:
```
chmod +x /path/to/your/project/run.sh
```
---

## ✍️ Author
Developed by CoderRooster

---