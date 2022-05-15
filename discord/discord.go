package discord

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-monitorsv2/scrapers"
)

const authToken = ""

var siteSlice = make(map[string]*scrapers.Site)

func CreateDiscordClient() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + authToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Printf("%s: %s\n", m.Author, m.Content)
	// If the message is "ping" reply with "Pong!"
	if strings.Contains(m.Content, "$siteadd") {
		splitMessage := strings.Split(m.Content, " ")
		url := splitMessage[1]
		domain := strings.Split(url, "://")[1]
		channelID := splitMessage[len(splitMessage)-1][2:][:18]
		go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Starting Monitor for %s", domain))
		currentSite := scrapers.MonitorSite(url, s, channelID)
		siteSlice[currentSite.URL] = currentSite
		// siteSlice[currentSite.URL].IsStopped = true
	}

	if strings.Contains(m.Content, "$siteremove") {
		splitMessage := strings.Split(m.Content, " ")
		url := splitMessage[1]
		domain := strings.Split(url, "://")[1]
		go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Removing Monitor for %s", domain))
		isFound := false
		fmt.Println("Removing", domain)
		for k, v := range siteSlice {
			fmt.Println(v.URL, url)
			if v.URL == url {
				siteSlice[k].IsStopped = true
				isFound = true
			}
		}

		if !isFound {
			go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not present in monitor", domain))
		}
	}

	if strings.Contains(m.Content, "check") {
		for k, v := range siteSlice {
			fmt.Println(k, v.IsStopped)
		}
	}
}
