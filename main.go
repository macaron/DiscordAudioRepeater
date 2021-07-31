package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var (
	token string
	gID   string
	cID   string
	audioFilePath string
)

func main() {
	// signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	// config
	flag.StringVar(&token, "t", os.Getenv("DISCORD_TOKEN"), "Discord Authentication Token")
	flag.StringVar(&gID, "g", os.Getenv("DISCORD_GUILD_ID"), "Guild ID of the channel to join")
	flag.StringVar(&cID, "c", os.Getenv("DISCORD_CHANNEL_ID"), "Channel ID of the channel to join")
	flag.StringVar(&audioFilePath, "p", os.Getenv("AUDIO_FILE_PATH"), "Specific Audio file path or URL")
	flag.Parse()
	if isEmptyFlag() {
		os.Exit(1)
	}

	d, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("[ERROR] failed to create session: ", err)
		return
	}
	defer d.Close()
	d.Identify.Intents = discordgo.IntentsGuildVoiceStates

	err = d.Open()
	if err != nil {
		fmt.Println("[ERROR] failed to open connection: ", err)
		return
	}

	v, err := d.ChannelVoiceJoin(gID, cID, false, false)
	if err != nil {
		fmt.Println("failed to join voice channel: ", err)
		return
	}
	fmt.Println("[INFO] joined voice channel")

	stop := make(chan bool)
	go func() {
		sig := <-sigs
		fmt.Printf("SIGNAL %d received, then application exit\n", sig)
		stop <- true
	}()

	fmt.Println("[INFO] playing audio")
	dgvoice.PlayAudioFile(v, audioFilePath, stop)
}

func isEmptyFlag() bool {
	if token == "" {
		fmt.Println("[WARN] Authentication Token is not defined")
	}
	if gID == "" {
		fmt.Println("[WARN] Guild ID is not defined")
	}
	if cID == "" {
		fmt.Println("[WARN] ChannelID is not defined")
	}
	if token == "" || gID == "" || cID == "" {
		return true
	}
	return false
}
