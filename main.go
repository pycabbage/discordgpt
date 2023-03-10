package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/comail/colog"
	"github.com/joho/godotenv"
	"github.com/pycabbage/discordgpt/internal/gpt"
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("warn: loading .env file")
	}
}

// func main() {
// 	res, err := gpt.CreateChatCompletion(context.Background(), gpt.ChatCompletionReq{
// 		Model: "gpt-3.5-turbo",
// 		Messages: []gpt.Message{
// 			{
// 				Role: gpt.MessageRoleSystem,
// 				Content: `
// 				あなたにはDiscord内のChatbotとしてユーザーと会話をしてもらいます。
// 				以下の制約条件を厳密に守って会話を行ってください。

// 				- セクシャルな話題に関しては誤魔化してください
// 				- なるべくシンプルな会話を心がけてください
// 				`,
// 			},
// 			{
// 				Role:    gpt.MessageRoleUser,
// 				Content: "こんにちは",
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Print(err)
// 		return
// 	}
// 	fmt.Print(strings.TrimSpace(res.Choices[0].Message.Content))
// }

func main() {
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime,
	})
	colog.Register()

	// load .env
	Env_load()
	// Create a new session using the DISCORD_TOKEN environment variable from Railway
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Printf("error: Error while starting bot: %s", err)
		return
	}

	// Add the message handler
	dg.AddHandler(messageCreate)

	// Add the Guild Messages intent
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Connect to the gateway
	err = dg.Open()
	if err != nil {
		log.Printf("error: while connecting to gateway: %s", err)
		return
	}

	// Wait until Ctrl+C or another signal is received
	log.Printf("info: The bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Close the Discord session
	dg.Close()
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't proceed if the message author is a bot
	if m.Author.Bot {
		return
	}

	Mentioned := false
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			Mentioned = true
		}
	}
	if !Mentioned {
		return
	}

	discordChannelIDs := strings.Split(os.Getenv("DISCORD_CHANNELID"), " ")
	if contains(discordChannelIDs, m.ChannelID) {
		return
	}

	log.Printf("info: Message(%s): %s", m.ChannelID, m.Content)
	// return

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong 🏓")
		return
	}

	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "Choo choo! 🚅")
		return
	}

	replyTree, _ := getReplyTree(s, m.ChannelID, m.ID)
	messagesForGpt := []gpt.Message{
		{
			Role: gpt.MessageRoleSystem,
			Content: `
			あなたにはDiscord内のChatbotとしてユーザーと会話をしてもらいます。
			以下の制約条件を厳密に守って会話を行ってください。
			
			- セクシャルな話題に関しては誤魔化してください
			- なるべくシンプルな会話を心がけてください
			- 適宜、会話にジョークを交えてください
			`,
		},
	}
	for i, discordMessage := range replyTree {
		var role string
		log.Printf("info: Message Tree(%s): %s", discordMessage.ChannelID, discordMessage.Content)

		// ここもうちょっと綺麗に書きたいね〜
		if i%2 == 0 {
			role = gpt.MessageRoleUser
			if discordMessage.Author.Bot {
				s.ChannelMessageSend(m.ChannelID, "エラー: リプライは交互に行うようにしてください")
				return
			}
		} else {
			role = gpt.MessageRoleAssistant
			if !discordMessage.Author.Bot {
				s.ChannelMessageSend(m.ChannelID, "エラー: リプライは交互に行うようにしてください")
				return
			}
		}

		messagesForGpt = append(messagesForGpt, gpt.Message{
			Role:    role,
			Content: discordMessage.Content,
		})
	}

	res, err := gpt.CreateChatCompletion(context.Background(), gpt.ChatCompletionReq{
		Model:    "gpt-3.5-turbo",
		Messages: messagesForGpt,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("エラーが発生しました: %s", err.Error()))
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, strings.TrimSpace(res.Choices[0].Message.Content), m.Reference())
}

func getReplyTree(s *discordgo.Session, channelID, messageID string) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message
	for {
		m, err := s.ChannelMessage(channelID, messageID)
		if err != nil {
			return nil, err
		}
		messages = append([]*discordgo.Message{m}, messages...)
		if m.MessageReference == nil {
			break
		}
		channelID = m.MessageReference.ChannelID
		messageID = m.MessageReference.MessageID
	}
	return messages, nil
}
