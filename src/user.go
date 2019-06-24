package src

const (
	ApiId    = 661076
	ApiHash  = "******************"
	ChatName = "CryptoMoodOfficial"
)

const (
	FilenameActiveUsers    = "activeUsers.csv"
	FilenameNotActiveUsers = "notActiveUsers.csv"
	FilenameBotUsers       = "botUsers.csv"
	FilenameMembers        = "mapMembersInfo.csv"
)

type User struct {
	Id int32 `json:"id"`
	// First name of the user
	FirstName string `json:"first_name"`
	// Last name of the user
	LastName string `json:"last_name"`
	// Username of the user
	Username string `json:"username"`
	// Phone number of the user
	PhoneNumber string `json:"phone_number"`
	// Current online status of the user
	Status string
	//The last visit time
	Timestamp int32
	//The date when user joined to the chat
	JoinChatDate int32
}


