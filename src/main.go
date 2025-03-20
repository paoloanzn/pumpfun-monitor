package main

import (
    "fmt"
    "os"
    "sync"
)

func startMonitorWithReconnection(f func(), maxRecon int, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < maxRecon; i++ {
        f()
    }
}

func main() {
	logger := NewLogger()
    var wg sync.WaitGroup

	if len(os.Args) < 2 {
        fmt.Println("Usage: pumpfun-monitor <command>")
        os.Exit(1)
    }

	switch os.Args[1] {
    case "start":
		mintMessageQueue, err := createNewMessageQueue()
		if err != nil {
			logger.Error("Error:", err)
		}
        migrationMessageQueue, err := createNewMessageQueue()
		if err != nil {
			logger.Error("Error:", err)
		}

        mintMonitor := func() {
		    startNewTokenMintMonitor(mintMessageQueue, logger)
        }

        migrationMonitor := func() {
            startNewMigrationMonitor(migrationMessageQueue, logger)
        }

        wg.Add(2)

        go startMonitorWithReconnection(mintMonitor, 100, &wg)
        go startMonitorWithReconnection(migrationMonitor, 100, &wg)

        wg.Wait()

    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}