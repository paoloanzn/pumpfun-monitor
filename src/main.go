package main

import (
    "fmt"
    "os"
    "sync"
    "time"
    "flag"
    "encoding/json"
)

type Config struct {
	MaxReconnections  int
	MintWorkers       int
	MigrationWorkers  int
}

func printUsage() {
	fmt.Printf(`Usage: %s <command> [options]

Commands:
  start    Start monitoring service
  version  Show version information

Start Options:
  -max-recon           Maximum reconnection attempts (default: 100)
  -mint-workers         Number of mint monitoring workers (default: 1)
  -migration-workers    Number of migration monitoring workers (default: 1)
`, os.Args[0])
	
	cmd := flag.NewFlagSet("", flag.ExitOnError)
	cmd.PrintDefaults()
}

func parseStartFlags() *Config {
	cmd := flag.NewFlagSet("start", flag.ExitOnError)
	config := &Config{
		// Default values
		MaxReconnections: 100,
		MintWorkers:      1,
		MigrationWorkers: 1,
	}

	cmd.IntVar(&config.MaxReconnections, "max-recon", config.MaxReconnections, 
		"Maximum reconnection attempts")
	cmd.IntVar(&config.MintWorkers, "mint-workers", config.MintWorkers, 
		"Number of mint monitoring workers")
	cmd.IntVar(&config.MigrationWorkers, "migration-workers", config.MigrationWorkers, 
		"Number of migration monitoring workers")

	cmd.Parse(os.Args[2:])
	return config
}



func startMonitorWithReconnection(f func(), maxRecon int, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < maxRecon; i++ {
        f()
        time.Sleep(time.Duration(i) * time.Second) 
    }
}

func createWebSocketConsumers(uuids []string, queue *MessageQueue) error {
    for _, uuid := range uuids {
        queue.RegisterConsumer(uuid)
        err := createWebSocketServer(uuid, queue)
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
	logger := NewLogger()
    var wg sync.WaitGroup

	if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

	switch os.Args[1] {
    case "start":
        config := parseStartFlags()

        mintMessageQueue, err := createNewMessageQueue()
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        migrationMessageQueue, err := createNewMessageQueue()
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        mintMonitor := func() {
		    startNewTokenMintMonitor(mintMessageQueue, logger)
        }

        migrationMonitor := func() {
            startNewMigrationMonitor(migrationMessageQueue, logger)
        }

        wg.Add(2)

        go startMonitorWithReconnection(mintMonitor, config.MaxReconnections, &wg)
        go startMonitorWithReconnection(migrationMonitor, config.MaxReconnections, &wg)

        mintConsumersIds, err := generateUUIDs(config.MintWorkers)
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        migrationConsumersIds, err := generateUUIDs(config.MigrationWorkers)
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        err = createWebSocketConsumers(mintConsumersIds, mintMessageQueue) 
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        err = createWebSocketConsumers(migrationConsumersIds, migrationMessageQueue) 
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        err = startWebSocketServers(8080)
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        mintConsumers := map[string]interface{}{
            "uuids": mintConsumersIds,
        }

        migrationConsumers := map[string]interface{}{
            "uuids": migrationConsumersIds,
        }

        mintConsumersJsonData, err := json.MarshalIndent(mintConsumers, "", "    ")
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        migrationConsumersJsonData, err := json.MarshalIndent(migrationConsumers, "", "    ")
        if err != nil {
            logger.Error("Error:", err)
            os.Exit(1)
        }

        logger.Info("Mint Consumers: %s", string(mintConsumersJsonData))
        logger.Info("Migration Consumers: %s", string(migrationConsumersJsonData))

        wg.Wait()

    case "version":
        fmt.Printf("%s\n", VERSION)

    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}