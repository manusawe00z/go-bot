package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-bot/internal/bot"
	"go-bot/internal/logging"
)

func main() {
	// Create a new bot instance
	b, err := bot.NewBot()
	if err != nil {
		fmt.Println("Error creating bot:", err)
		os.Exit(1)
	}

	// Start the bot
	if err := b.Start(); err != nil {
		fmt.Println("Error starting bot:", err)
		os.Exit(1)
	}
	defer b.Stop()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait for interrupt signal to gracefully shut down
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-sc

	logging.Info("Received signal %v, shutting down bot...", sig)
}
