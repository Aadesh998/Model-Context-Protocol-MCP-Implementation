package logger

import (
	"log"
	"os"
)

var Log *log.Logger

func Init() {
	file, err := os.OpenFile("mcp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}

	Log = log.New(file, "MCP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
