package scrapers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Site struct {
	URL            string
	DiscordSession *discordgo.Session
	ChannelID      string
	IsStopped      bool
}

// type Site struct {
// 	Site types.Site
// }

func MonitorSite(url string, s *discordgo.Session, channelID string) *Site {
	m := &Site{}
	m.IsStopped = false
	m.URL = url
	m.DiscordSession = s
	m.ChannelID = channelID
	fmt.Println(m.URL)
	go m.initMonitor()
	return m
}

func (m *Site) initMonitor() {
	var currentBody []byte = []byte("Test")
	for !m.IsStopped {
		fmt.Println("Checking Site ", m.URL, m.IsStopped)
		body, err := checkSite(m.URL)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if bytes.Equal(currentBody, []byte("Test")) {
			currentBody = body
			continue
		}

		if bytes.Compare(body, currentBody) != 0 {
			currentBody = body
			fmt.Println("Difference Detected")
			title := strings.Split(strings.Split(string(currentBody), "<title>")[1], "</title>")[0]
			desc := strings.Split(strings.Split(string(currentBody), `name="description" content="`)[1], ">")[0]
			img := strings.Split(strings.Split(string(currentBody), `<meta property="og:image" content="`)[1], `"`)[0]
			// image := &discordgo.MessageEmbedImage{URL: img}
			embed := &discordgo.MessageEmbed{
				URL:  m.URL,
				Type: "link",
				Description: "**Description :** \n" + desc,
				// Image: image,
				Title: title,
				Author: &discordgo.MessageEmbedAuthor{Name: "Price Errors Site Change", IconURL: "https://cdn.discordapp.com/attachments/972640827396477048/974416841663475712/image1.png"},
				Thumbnail: &discordgo.MessageEmbedThumbnail{URL: img},
			}
			fmt.Println(desc, "\n\n\n", img, embed)
			// m.DiscordSession.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", title))
			fmt.Println(m.DiscordSession.ChannelMessageSendEmbed(m.ChannelID, embed))
		}

		time.Sleep(5 * time.Second)

	}
	return
}

func (m *Site) Stop() {
	m.IsStopped = true
}
func checkSite(url string) ([]byte, error) {

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return []byte(""), err
	}
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	// req.Header.Add("accept", "application/json")
	req.Header.Add("dnt", "1")
	req.Header.Add("accept-language", "en")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	// req.Header.Set("Connection", "close")
	req.Close = true
	// Fetch Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return []byte(""), err
	}
	defer resp.Body.Close()
	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return []byte(""), err
	}
	return respBody, nil

	// // Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))
}
