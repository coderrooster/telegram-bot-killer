package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// example: https://api.telegram.org/bot1234567890:AAG-abcdefghijklmnopqrstuvwxyz/sendMessage?parse_mode=markdown&chat_id=1234567890&text="test message"

func sendMessage(i int, logger *log.Logger) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TOKEN"))

	body := map[string]interface{}{
		"chat_id": os.Getenv("CHAT_ID"),
		"text":    fmt.Sprintf("%s #%d", os.Getenv("CHAT_MESSAGE"), i),
	}
	bodyJSON, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		logger.Printf("[ERROR] Failed to send request #%d: %v\n", i, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Baca isi respons error dari Telegram
		var resBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&resBody)

		logger.Printf("[ERROR] Send #%d failed - Status: %d - Response: %s\n", i, resp.StatusCode, resBody)
		if resp.StatusCode == 429 {
			return fmt.Errorf("Too Many Requests")
		}
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	logger.Printf("✅ Send #%d - Status: %s\n", i, resp.Status)
	return nil
}

func main() {
	// get current working directory
	folderPath, err := os.Getwd()
	if err != nil {
		fmt.Println("[FATAL] Failed to get working directory:", err)
		return
	}

	// file name with format "log-YYYY-MM-DD.log"
	dateStr := time.Now().Format("2006-01-02")
	logFilename := fmt.Sprintf("log-%s.log", dateStr)
	logPath := filepath.Join(folderPath, logFilename)

	// open or create log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("[FATAL] Failed to open log file:", err)
		return
	}
	defer file.Close()

	// Logger ke file + console
	multiWriter := io.MultiWriter(file, os.Stdout)
	logger := log.New(multiWriter, "", log.LstdFlags)

	logger.Println("=== Starting Telegram Bot Kill ===")

	// set request tries from env variable, default to 100
	requestTriesStr := os.Getenv("REQUEST_EVERY_TRIES")
	requestTries := 100 // default value
	if requestTriesStr != "" {
		if v, err := strconv.Atoi(requestTriesStr); err == nil {
			requestTries = v
		} else {
			logger.Printf("[WARN] Invalid REQUEST_EVERY_TRIES value: %v. Using default 100.\n", err)
		}
	}

	for i := 1; i <= requestTries; i++ {
		err := sendMessage(i, logger)
		if err != nil {
			logger.Println("[INFO] Error occurred, waiting 5 seconds before retrying...")
			logger.Printf("[INFO] Error: %v\n", err)

			sleepTimeStr := os.Getenv("SLEEP_TIME")
			sleepTime := 5 // default to 5 seconds if not set or invalid
			if sleepTimeStr != "" {
				if v, err := strconv.Atoi(sleepTimeStr); err == nil {
					sleepTime = v
				} else {
					logger.Printf("[WARN] Invalid SLEEP_TIME value: %v. Using default 5 seconds.\n", err)
				}
			}
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}

	logger.Println("=== Telegram Bot Kill Finished ===")
	if err := file.Sync(); err != nil {
		logger.Printf("[ERROR] Failed to sync log file: %v\n", err)
	} else {
		logger.Println("Log file synced successfully.")
	}
}
