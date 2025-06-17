package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	cachedMessage       string
	lastGeneratedMinute int
)

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getCachedOrNewMessage(length int) string {
	currentMinute := time.Now().Minute()
	if currentMinute != lastGeneratedMinute || cachedMessage == "" {
		cachedMessage = generateRandomString(length)
		lastGeneratedMinute = currentMinute
	}
	return cachedMessage
}

func getEnvString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func sendMessage(i int, logger *log.Logger) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", getEnvString("BOT_TOKEN", ""))

	messageLength := getEnvInt("CHAT_MESSAGE_LENGTH", 1000)
	if messageLength > 4000 {
		messageLength = 4000
	}

	chatText := getCachedOrNewMessage(messageLength)

	body := map[string]interface{}{
		"chat_id": getEnvString("CHAT_ID", ""),
		"text":    fmt.Sprintf("%s", chatText),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		logger.Printf("[ERROR] Failed to marshal request body: %v\n", err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		logger.Printf("[ERROR] Failed to send request #%d: %v\n", i, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var resBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
			logger.Printf("[ERROR] Failed to parse error response: %v\n", err)
		}
		logger.Printf("[ERROR] Send #%d failed - Status: %d - Response: %v\n", i, resp.StatusCode, resBody)

		if resp.StatusCode == 429 {
			return fmt.Errorf("Too Many Requests")
		}
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	logger.Printf("✅ Send #%d - Status: %s\n", i, resp.Status)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	folderPath, err := os.Getwd()
	if err != nil {
		fmt.Println("[FATAL] Failed to get working directory:", err)
		return
	}

	dateStr := time.Now().Format("2006-01-02")
	logFilename := fmt.Sprintf("log-%s.log", dateStr)
	logPath := filepath.Join(folderPath, logFilename)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("[FATAL] Failed to open log file:", err)
		return
	}
	defer file.Close()

	multiWriter := io.MultiWriter(file, os.Stdout)
	logger := log.New(multiWriter, "", log.LstdFlags)

	if getEnvString("BOT_TOKEN", "") == "" || getEnvString("CHAT_ID", "") == "" {
		logger.Println("[FATAL] BOT_TOKEN or CHAT_ID is not set.")
		return
	}

	logger.Println("=== Starting Telegram Bot Kill ===")

	requestTries := getEnvInt("REQUEST_EVERY_TRIES", 100)
	sleepTime := getEnvInt("SLEEP_TIME", 5)

	for i := 1; i <= requestTries; i++ {
		err := sendMessage(i, logger)
		if err != nil {
			logger.Println("[INFO] Error occurred, waiting before retrying...")
			logger.Printf("[INFO] Error: %v\n", err)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}

	logger.Println("=== Telegram Bot Kill Finished ===")
	if err := file.Sync(); err != nil {
		logger.Printf("[ERROR] Failed to sync log file: %v\n", err)
	} else {
		logger.Println("Log file synced successfully.")
	}
}
