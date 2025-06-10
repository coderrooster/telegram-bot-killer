#!/bin/bash

# === BASE CONFIG ===
PROJECT_DIR="/www/wwwroot/telegram-bot-killer"
ENV_FILE="$PROJECT_DIR/.env"
LOG_DIR="$PROJECT_DIR/logs-cron"
LOG_FILE="$LOG_DIR/log-cron-$(date +'%Y-%m-%d').log"
GO_FILE="$PROJECT_DIR/kill-bot.go"
BINARY="$PROJECT_DIR/kill-bot"
RETENTION_DAYS=30

# === enter your project directory here ===
cd "$PROJECT_DIR" || {
    echo "[FATAL] Failed to cd to $PROJECT_DIR"
    exit 1
}

# === make directory for logs if not exists ===
mkdir -p "$LOG_DIR"

# === Load environment from .env file ===
if [ ! -f "$ENV_FILE" ]; then
    echo "[FATAL] .env file not found: $ENV_FILE" >> "$LOG_FILE"
    exit 1
fi

# load all environment variables from .env file except comments and empty lines
export $(grep -v '^#' "$ENV_FILE" | xargs)

# === Build binary if not exists ===
if [ ! -f "$BINARY" ]; then
    echo "[INFO] Binary not found. Building from $GO_FILE ..." >> "$LOG_FILE"
    go build -o "$BINARY" "$GO_FILE" >> "$LOG_FILE" 2>&1

    if [ $? -ne 0 ]; then
        echo "[FATAL] Build failed. Check Go source." >> "$LOG_FILE"
        exit 1
    fi

    echo "[INFO] Build completed." >> "$LOG_FILE"
fi

# === run the binary and log output ===
echo "[INFO] === $(date '+%Y-%m-%d %H:%M:%S') - Run Started ===" >> "$LOG_FILE"
"$BINARY" >> "$LOG_FILE" 2>&1
echo "[INFO] === $(date '+%Y-%m-%d %H:%M:%S') - Run Finished ===" >> "$LOG_FILE"

# === delete logs older than RETENTION_DAYS ===
find "$LOG_DIR" -name "log-cron-*.log" -type f -mtime +$RETENTION_DAYS -exec rm {} \;
