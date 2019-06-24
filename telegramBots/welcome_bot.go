package telegramBots

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"telegramBot/src"
)


func ProcessWelcomeBot() {

	tdlibClient := src.Init()
	listener := tdlibClient.GetListener()
	defer listener.Close()

	var memberCounter int32
	chatInfoRequest := &client.SearchPublicChatRequest{Username: src.ChatName}
	chatInfo, err := tdlibClient.SearchPublicChat(chatInfoRequest)
	if err != nil {
		log.Fatalf("Search Public Info: %s", err)
	}

	var oldMemberCounter = src.GetSuperGroupInfo(*tdlibClient, chatInfo.Id).MemberCount
	for update := range listener.Updates {

		if update.GetClass() == client.ClassUpdate && update.GetType() == client.TypeUpdateNewMessage {
			message := update.(*client.UpdateNewMessage).Message


			log.Println(chatInfo)
			//If it is in selected chat and it is a system msg
			if message.ChatId == chatInfo.Id && !message.CanBeForwarded {

				memberCounter = src.GetSuperGroupInfo(*tdlibClient, chatInfo.Id).MemberCount
				//Check if someone join to group
				if memberCounter > oldMemberCounter {

					userRequest := &client.GetUserRequest{UserId: message.SenderUserId}
					user, err := tdlibClient.GetUser(userRequest)
					if err != nil {
						log.Fatalf("GetUserFullInfo: %s", err)
					}
					var setWelcomeMsg *client.InputMessageText
					if user.Username == "" {
						if user.LastName == "" {
							setWelcomeMsg = &client.InputMessageText{Text: &client.FormattedText{Text: user.FirstName+" welcome to the CryptoMood telegram group!🖖 CryptoMood is a cutting-edge app that uses the power of AI & data from 50,000 sources to improve your trading/investing decisions."}}
						} else {
							setWelcomeMsg = &client.InputMessageText{Text: &client.FormattedText{Text: user.FirstName+" "+user.LastName+" welcome to the CryptoMood telegram group!🖖 CryptoMood is a cutting-edge app that uses the power of AI & data from 50,000 sources to improve your trading/investing decisions."}}
						}

					} else {
						setWelcomeMsg = &client.InputMessageText{Text: &client.FormattedText{Text: "@" + user.Username + " welcome to the CryptoMood telegram group!🖖 CryptoMood is a cutting-edge app that uses the power of AI & data from 50,000 sources to improve your trading/investing decisions."}}
					}
					welcomeMsg := client.InputMessageContent(setWelcomeMsg)

					welcomeMsgRequest := &client.SendMessageRequest{ChatId: chatInfo.Id, InputMessageContent: welcomeMsg}
					_, err = tdlibClient.SendMessage(welcomeMsgRequest)
					if err != nil {
						log.Fatalf("Send Message: %s", err)
					}

				}

				oldMemberCounter = memberCounter
			}

		}

	}
}
