package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Constants for commands and statuses
const (
	AddCommand      = 0
	DeleteCommand   = 1
	ContinueCommand = 2
	AbortCommand    = 3
	OSSuccess       = 0
	OSInvalid       = -1
	LogFileWindows  = `C:\Program Files (x86)\ossec-agent\active-response\active-responses.log`
	LogFileLinux    = "/var/ossec/logs/active-responses.log"
)

// Message represents the structure for input data
type Message struct {
	Alert   map[string]interface{} `json:"alert"`
	Command int                    `json:"command"`
}

// getLogFilePath returns the appropriate log file path based on the OS
func getLogFilePath() string {
	if os.PathSeparator == '\\' {
		return LogFileWindows
	}
	return LogFileLinux
}

// writeDebugFile writes messages to the debug log
func writeDebugFile(arName, msg string) error {
	logFile := getLogFilePath()
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	arNamePosix := strings.ReplaceAll(arName, "\\", "/")
	entry := fmt.Sprintf("%s %s: %s\n", time.Now().Format("2006/01/02 15:04:05"), arNamePosix, msg)
	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}
	return nil
}

// setupAndCheckMessage parses and validates input from stdin
func setupAndCheckMessage(arName string) (Message, error) {
	reader := bufio.NewReader(os.Stdin)
	inputStr, err := reader.ReadString('\n')
	if err != nil {
		writeDebugFile(arName, "Failed to read stdin")
		return Message{}, fmt.Errorf("failed to read stdin: %w", err)
	}

	writeDebugFile(arName, inputStr)

	var msg Message
	if err := json.Unmarshal([]byte(inputStr), &msg.Alert); err != nil {
		writeDebugFile(arName, "Decoding JSON has failed, invalid input format")
		return Message{}, fmt.Errorf("invalid input format: %w", err)
	}

	command, ok := msg.Alert["command"].(string)
	if !ok {
		writeDebugFile(arName, "Missing or invalid 'command' field")
		return Message{}, errors.New("missing or invalid 'command' field")
	}

	switch command {
	case "add":
		msg.Command = AddCommand
	case "delete":
		msg.Command = DeleteCommand
	default:
		msg.Command = OSInvalid
		writeDebugFile(arName, "Invalid command: "+command)
	}

	return msg, nil
}

// sendKeysAndCheckMessage sends keys and validates the response
func sendKeysAndCheckMessage(arName string, keys []interface{}) (int, error) {
	message := map[string]interface{}{
		"version": 1,
		"origin": map[string]string{
			"name":   filepath.Base(arName),
			"module": "active-response",
		},
		"command":    "check_keys",
		"parameters": map[string]interface{}{"keys": keys},
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		writeDebugFile(arName, "Failed to encode message to JSON")
		return OSInvalid, fmt.Errorf("failed to encode message: %w", err)
	}

	writeDebugFile(arName, string(messageJSON))
	fmt.Println(string(messageJSON))

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		writeDebugFile(arName, "Failed to read response")
		return OSInvalid, fmt.Errorf("failed to read response: %w", err)
	}

	writeDebugFile(arName, response)

	var responseMsg map[string]interface{}
	if err := json.Unmarshal([]byte(response), &responseMsg); err != nil {
		writeDebugFile(arName, "Decoding JSON has failed, invalid input format")
		return OSInvalid, fmt.Errorf("invalid response format: %w", err)
	}

	action, ok := responseMsg["command"].(string)
	if !ok {
		writeDebugFile(arName, "Invalid or missing 'command' field in response")
		return OSInvalid, errors.New("invalid or missing 'command' field")
	}

	switch action {
	case "continue":
		return ContinueCommand, nil
	case "abort":
		return AbortCommand, nil
	default:
		writeDebugFile(arName, "Invalid 'command' value in response")
		return OSInvalid, errors.New("invalid 'command' value in response")
	}
}

// main function
func main() {
	arName := os.Args[0]

	if err := writeDebugFile(arName, "Started"); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to debug log:", err)
		os.Exit(OSInvalid)
	}

	msg, err := setupAndCheckMessage(arName)
	if err != nil || msg.Command < 0 {
		os.Exit(OSInvalid)
	}

	switch msg.Command {
	case AddCommand:
		alert, ok := msg.Alert["parameters"].(map[string]interface{})["alert"].(map[string]interface{})
		if !ok {
			writeDebugFile(arName, "Invalid alert format")
			os.Exit(OSInvalid)
		}

		keys := []interface{}{alert["rule"].(map[string]interface{})["id"]}
		action, err := sendKeysAndCheckMessage(arName, keys)
		if err != nil || action != ContinueCommand {
			if action == AbortCommand {
				writeDebugFile(arName, "Aborted")
				os.Exit(OSSuccess)
			}
			writeDebugFile(arName, "Invalid command")
			os.Exit(OSInvalid)
		}

		file, err := os.OpenFile("ar-test-result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			writeDebugFile(arName, "Failed to write result file")
			os.Exit(OSInvalid)
		}
		defer file.Close()

		if _, err := file.WriteString(fmt.Sprintf("Active response triggered by rule ID: <%v>\n", keys)); err != nil {
			writeDebugFile(arName, "Failed to write result file")
			os.Exit(OSInvalid)
		}

	case DeleteCommand:
		if err := os.Remove("ar-test-result.txt"); err != nil {
			writeDebugFile(arName, "Failed to delete result file")
			os.Exit(OSInvalid)
		}

	default:
		writeDebugFile(arName, "Invalid command")
	}

	writeDebugFile(arName, "Ended")
	os.Exit(OSSuccess)
}
