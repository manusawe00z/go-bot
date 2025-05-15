package main

import (
	"fmt"
	"os"
	"os/signal"

	"gobot/internal/bot"
)

func main() {
	b, err := bot.NewBot()
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	err = b.Start()
	if err != nil {
		fmt.Println("Error starting bot:", err)
		return
	}
	defer b.Stop()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// รอจนกว่าจะได้รับสัญญาณหยุด
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Shutting down bot.")
}
