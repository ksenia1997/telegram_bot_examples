package telegramBots

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"telegramBot/src"
	"time"
)

func getMsgWithOldestDate(data map[int64]src.Message) src.Message {
	var oldestMsg src.Message
	var counter int64
	for _, msg := range data {
		if counter == 0 {
			oldestMsg = msg
			counter++
			continue
		}
		if time.Unix(int64(msg.Date), 0).Sub(time.Unix(int64(oldestMsg.Date), 0)) < 0 {
			oldestMsg = msg
		}
	}

	return oldestMsg
}
func ProcessInvitation() {

	tdlibClient := src.Init()
	listener := tdlibClient.GetListener()
	defer listener.Close()

	//get chat id int32
	chatInfoRequest := &client.SearchPublicChatRequest{Username: src.ChatName}
	chatInfo, err := tdlibClient.SearchPublicChat(chatInfoRequest)
	if err != nil {
		log.Fatalf("Search Public Info: %s", err)
	}

	var msgId int64
	var msgMap = src.LoadMessagesFromCSV(src.FilenameMessages)
	if msgMap == nil {
		msgMap = make(map[int64]src.Message)
	} else {
		msgId = getMsgWithOldestDate(msgMap).MsgId
	}

	chatHistoryRequest := &client.GetChatHistoryRequest{ChatId: chatInfo.Id, FromMessageId: msgId, Offset: -99, Limit: 100, OnlyLocal: false}
	chatHistory, err := tdlibClient.GetChatHistory(chatHistoryRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range chatHistory.Messages {

		msgId = message.Id
		var contentText string
		var contentType = message.Content.MessageContentType()
		if contentType == "messageText" {
			contentText = message.Content.(*client.MessageText).Text.Text
		} else if contentType == "messageChatAddMembers" {
			continue
		}
		msgMap[msgId] = src.Message{MsgId: msgId, ContentType: contentType, ContentText: contentText, SenderUserId: message.SenderUserId, Date: message.Date, ReplyToMessageId: message.ReplyToMessageId}

	}

	for {
		chatHistoryRequest = &client.GetChatHistoryRequest{ChatId: chatInfo.Id, FromMessageId: msgId, Offset: -99, Limit: 100, OnlyLocal: false}
		chatHistory, err = tdlibClient.GetChatHistory(chatHistoryRequest)
		if err != nil {
			src.SaveMsgToCSV(msgMap, src.FilenameMessages)
			log.Fatal(err)
		}
		for _, msg := range chatHistory.Messages {
			msgId = msg.Id
			var contentText string
			var contentType = msg.Content.MessageContentType()
			if contentType == "messageText" {
				contentText = msg.Content.(*client.MessageText).Text.Text
			} else if contentType == "messageChatAddMembers" {
				continue
			}
			msgMap[msgId] = src.Message{MsgId: msgId, ContentType: contentType, ContentText: contentText, SenderUserId: msg.SenderUserId, Date: msg.Date, ReplyToMessageId: msg.ReplyToMessageId}
		}
		log.Println(len(msgMap))
	}

}
