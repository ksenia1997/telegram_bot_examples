package src

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"path/filepath"
)

func Init() *client.Client {
	client.SetLogVerbosityLevel(1)

	// client authorizer
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  ApiId,
		ApiHash:                ApiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	me, err := tdlibClient.GetMe()
	if err != nil {
		log.Fatalf("GetMe error: %s", err)
	}

	log.Printf("Me: %s %s [%s]", me.FirstName, me.LastName, me.Username)
	return tdlibClient
}


func GetSuperGroupId(tdlibClient *client.Client, id int64) int32 {
	chatRequest := &client.GetChatRequest{ChatId: id}
	chat, err := tdlibClient.GetChat(chatRequest)
	if err != nil {
		log.Fatalf("GetChat error: %s", err)
	}
	sup := chat.Type.(*client.ChatTypeSupergroup)
	return sup.SupergroupId
}



func GetSuperGroupInfo(tdlibClient client.Client, id64 int64) *client.SupergroupFullInfo {
	supRequest := &client.GetSupergroupFullInfoRequest{SupergroupId: GetSuperGroupId(&tdlibClient, id64)}
	superGroup, err := tdlibClient.GetSupergroupFullInfo(supRequest)
	if err != nil {
		log.Fatalf("GetSupergroupFullInfo error: %s", err)
	}
	return superGroup
}

