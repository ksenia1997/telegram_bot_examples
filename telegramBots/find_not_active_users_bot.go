package telegramBots

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"os"
	"strconv"
	"telegramBot/src"
	"time"
)

const (
	layout     = "2006-01-02T15:04:05.000Z"
	acceptTime = "2019-05-01T00:00:00.000Z"
)

func getMembers(isFirstRequest bool, groupInfo *client.SupergroupFullInfo, superGroupId int32, tdlibClient *client.Client) {
	var mapMembers map[int32]src.User
	if isFirstRequest {
		mapMembers = make(map[int32]src.User)
	} else {
		mapMembers = src.LoadUsersInfoFromCSV(src.FilenameMembers)
	}
	var offset int32
	var nonBotCounter = 0
	log.Println(len(mapMembers))
	for offset = int32(len(mapMembers)); offset < groupInfo.MemberCount; offset += 200 {
		groupMembersRequest := &client.GetSupergroupMembersRequest{SupergroupId: superGroupId, Offset: offset, Limit: 5000}
		groupMembers, err := tdlibClient.GetSupergroupMembers(groupMembersRequest)
		if err != nil {
			src.SaveToCSV(mapMembers, src.FilenameMembers)
			log.Fatalf("GetSupergroupMembers: %s", err)
		}

		for _, member := range groupMembers.Members {
			if member.BotInfo == nil {
				nonBotCounter++
				if _, ok := mapMembers[member.UserId]; !ok {
					mapMembers[member.UserId] = src.User{Id: member.UserId, JoinChatDate: member.JoinedChatDate}
				}

			}
		}

		log.Println("NonBot: ", nonBotCounter)
		log.Println("Len of Map: ", len(mapMembers))

	}
	src.SaveToCSV(mapMembers, src.FilenameMembers)

}

func detectUnusedAccounts(tdlibClient *client.Client) {

	members := src.LoadUsersInfoFromCSV(src.FilenameMembers)

	acceptedTimestamp, err := time.Parse(layout, acceptTime)
	if err != nil {
		log.Println(err)
	}

	var notActiveUsers map[int32]src.User
	var botUsers map[int32]src.User
	var activeUsers map[int32]src.User

	notActiveUsers = src.LoadUsersInfoFromCSV(src.FilenameNotActiveUsers)
	if notActiveUsers == nil {
		notActiveUsers = make(map[int32]src.User)
	}

	botUsers = src.LoadUsersInfoFromCSV(src.FilenameBotUsers)
	if botUsers == nil {
		botUsers = make(map[int32]src.User)
	}

	activeUsers = src.LoadUsersInfoFromCSV(src.FilenameActiveUsers)
	if activeUsers == nil {
		activeUsers = make(map[int32]src.User)
	}

	log.Println("Len of active users: ", len(activeUsers))
	log.Println("Len of non active users: ", len(notActiveUsers))
	log.Println("Len of bots: ", len(botUsers))

	for memberId, memberInfo := range members {

		_, isNotActiveUser := notActiveUsers[memberId]
		_, isActiveUser := activeUsers[memberId]
		_, isBotUser := botUsers[memberId]
		if isNotActiveUser || isActiveUser || isBotUser {
			continue
		}

		userRequest := &client.GetUserRequest{UserId: memberId}
		userInfo, err := tdlibClient.GetUser(userRequest)

		if err != nil {
			log.Println("active users: ", len(activeUsers))
			log.Println("not active Users: ", len(notActiveUsers))
			log.Println("bot Users: ", len(botUsers))
			src.SaveToCSV(notActiveUsers, src.FilenameNotActiveUsers)
			src.SaveToCSV(botUsers, src.FilenameBotUsers)
			src.SaveToCSV(activeUsers, src.FilenameActiveUsers)

			log.Fatalf("Get UserFullInfo: %s", err)
		}

		user := src.User{Id: userInfo.Id, FirstName: userInfo.FirstName, LastName: userInfo.LastName, Username: userInfo.Username, PhoneNumber: userInfo.PhoneNumber, Status: userInfo.Status.UserStatusType(), JoinChatDate: memberInfo.JoinChatDate}

		if userInfo.Status.UserStatusType() == "userStatusEmpty" {
			botUsers[memberId] = user
		} else if userInfo.Status.UserStatusType() == "userStatusLastMonth" || userInfo.Status.UserStatusType() == "userStatusLastWeek" {
			notActiveUsers[memberId] = user
		} else if userInfo.Status.UserStatusType() == "userStatusOffline" {

			user.Timestamp = userInfo.Status.(*client.UserStatusOffline).WasOnline


			if int64(userInfo.Status.(*client.UserStatusOffline).WasOnline) < acceptedTimestamp.Unix() {
				if (userInfo.Status.(*client.UserStatusOffline).WasOnline - memberInfo.JoinChatDate) < 259200 {
					botUsers[memberId] = user
				} else {
					notActiveUsers[memberId] = user
				}
			} else {
				activeUsers[memberId] = user
			}
		} else {
			activeUsers[memberId] = user
		}

	}
	log.Println("active users: ", len(activeUsers))
	log.Println("not active Users: ", len(notActiveUsers))
	log.Println("bot Users: ", len(botUsers))
	src.SaveToCSV(notActiveUsers, src.FilenameNotActiveUsers)
	src.SaveToCSV(botUsers, src.FilenameBotUsers)
	src.SaveToCSV(activeUsers, src.FilenameActiveUsers)
}

func ProcessFindUsers() {

	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatalf("\nPlease enter as a first argument - \"members\" or \"usersAccounts\" \nAs a second argument if it is a first start of this program or not - \"true\" or \"false\". \n")
	}

	if args[0] != "members" && args[0] != "userAccounts" {
		log.Fatal("The first argument must be \"members\" or \"usersAccounts\" \n")
	}

	if args[1] != "true" && args[1] != "false" {
		log.Fatal("The second argument must be \"true\" or \"false\" \n")
	}

	tdlibClient := src.Init()
	listener := tdlibClient.GetListener()
	defer listener.Close()

	//get chat id int32
	chatInfoRequest := &client.SearchPublicChatRequest{Username: src.ChatName}
	chatInfo, err := tdlibClient.SearchPublicChat(chatInfoRequest)
	if err != nil {
		log.Fatalf("Search Public Info: %s", err)
	}

	log.Println(chatInfo)
	superGroupId := src.GetSuperGroupId(tdlibClient, chatInfo.Id)

	// get info about group
	groupInfoRequest := &client.GetSupergroupFullInfoRequest{SupergroupId: superGroupId}
	groupInfo, err := tdlibClient.GetSupergroupFullInfo(groupInfoRequest)
	if err != nil {
		log.Fatalf("GetSuperGroupFullInfo: %s", err)
	}
	log.Println(groupInfo)

	isFirstIteration, _ := strconv.ParseBool(args[1])
	if args[0] == "members" {
		getMembers(isFirstIteration, groupInfo, superGroupId, tdlibClient)
	} else {
		detectUnusedAccounts(tdlibClient)
	}

}
