package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	defer dg.Close()
	dg.AddHandler(onReady)
	dg.AddHandler(onMessage)
	must(err)
	must(dg.Open())
	log.Printf("Running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Printf("Received ABORT")
	dg.Close()

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {

}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	log.Printf("Message %s", m.Content)
	if strings.HasPrefix(m.Content, "!hello") {
		s.ChannelMessageSend(m.ChannelID, "Hello")
		log.Printf("Sent Message Hello")
	}
	if strings.HasPrefix(m.Content, "!poe") {
		r := strings.NewReplacer("!poe ", "")
		str := r.Replace(m.Content)
		blankViaUnderscore := strings.NewReplacer(" ", "_")
		urlFormatted := blankViaUnderscore.Replace(str)
		wikiurl := "https://pathofexile.gamepedia.com/" + urlFormatted

		resp, err := http.Get(wikiurl)
		must(err)
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			log.Printf("%s not found", wikiurl)
			s.ChannelMessageSend(m.ChannelID, "Not Found")
			return
		}

		s.ChannelMessageSend(m.ChannelID, wikiurl)

	}
}
