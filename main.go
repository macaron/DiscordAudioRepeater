package main

import (
	"bytes"
	json2 "encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	token          string
	gID            string
	cID            string
	audioFilePath  string
	discordWebhook string
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
	flag.StringVar(&discordWebhook, "w", os.Getenv("DISCORD_WEBHOOK_URL"), "Specific Discord webhook")
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
	throwWebhook("[INFO] joined voice channel")

	stop := make(chan bool)
	go func() {
		sig := <-sigs
		exitReason := fmt.Sprintf("[INFO] SIGNAL %d received, then application exit\n", sig)
		throwWebhook(exitReason)
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

func throwWebhook(content string) {
	fmt.Println(content)

	if discordWebhook == "" {
		return
	}

	type BodyJson struct {
		Content string `json:"content"`
	}
	json, _ := json2.Marshal(&BodyJson{Content: content})

	req, _ := http.NewRequest("POST", discordWebhook, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		fmt.Println("[ERROR] failed to send webhook")
	}
}
