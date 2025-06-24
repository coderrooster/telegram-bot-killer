#!/bin/bash

# === BASE CONFIG ===
PROJECT_DIR="/www/wwwroot/telegram-bot-killer"
ENV_FILE="$PROJECT_DIR/.env"
LOG_DIR="$PROJECT_DIR/logs-cron"
LOG_FILE="$LOG_DIR/log-$(date +'%Y-%m-%d').log"
GO_FILE="$PROJECT_DIR/kill-bot.go"
BINARY="$PROJECT_DIR/kill-bot"
DAYS_TO_KEEP=3
EXCLUDES=""

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
# export $(grep -v '^#' "$ENV_FILE" | xargs)

# Safe way to load .env with support for quotes and spaces
set -a
source "$ENV_FILE"
set +a

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

# === keep logs for a certain number of days ===
for ((i=0; i<DAYS_TO_KEEP; i++)); do
    EXCLUDES+="! -name \"log-$(date -d \"$i day ago\" +%Y-%m-%d).log\" "
done
echo "[INFO] Cleaning up old logs..." >> "$LOG_FILE"
eval "find \"$LOG_DIR\" -type f -name \"log-*.log\" $EXCLUDES -exec echo \"[INFO] Deleting: {}\" >> \"$LOG_FILE\" \; -exec rm {} \;"
eval "find \"$PROJECT_DIR\" -type f -name \"log-*.log\" $EXCLUDES -exec echo \"[INFO] Deleting: {}\" >> \"$LOG_FILE\" \; -exec rm {} \;"

# example of how to run the script:
# find "/www/wwwroot/telegram-bot-killer/logs-cron" -type f -name "log-*.log" \
# ! -name "log-2025-06-24.log" \
# ! -name "log-2025-06-23.log" \
# ! -name "log-2025-06-22.log" \
# -exec echo "Deleting: {}" \; \
# -exec rm {} \;

# === end of script ===
echo "[INFO] Script completed successfully." >> "$LOG_FILE"
exit 0