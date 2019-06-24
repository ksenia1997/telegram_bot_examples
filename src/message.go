package src

type Message struct {
	MsgId            int64
	ContentType      string
	ContentText      string
	SenderUserId     int32
	Date             int32
	ReplyToMessageId int64
}

type MessageWithUserInfo struct {
	MsgId       int64
	ContentType string
	ContentText string
	UserId      int32
	Username    string
}

const FilenameMessages = "msgsFromChat.csv"
