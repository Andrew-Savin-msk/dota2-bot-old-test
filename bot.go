package main

import (
	//"context"

	"fmt"
	f "fmt"
	"log"
	"reflect"
	"strings"

	"time"

	"encoding/json"
	"io/ioutil"

	"os"

	"github.com/paralin/go-dota2"
	"github.com/paralin/go-dota2/events"
	gcmm "github.com/paralin/go-dota2/protocol"
	"github.com/paralin/go-steam"
	"github.com/paralin/go-steam/protocol/steamlang"
	"github.com/paralin/go-steam/trade"
	"github.com/sirupsen/logrus"
)

const (
	botVersion string = "0.1.5"
	botBuild   string = "-dev"
)

var (
	connects = 1
)

func main() {
	info()

	readJSON()

	account := &steam.LogOnDetails{
		Username: "yandexProg",
		Password: "9@h$YHFAVLJHnUg",
	}
	// f.Println("Username:")
	// f.Scan(&account.Username)
	// f.Println("Password: ")
	// f.Scan(&account.Password)
	// f.Println("Steam Guard: ")
	// f.Scan(&account.TwoFactorCode)

	//Login with Hash | #InFuture
	hash, err := os.ReadFile("sentry")
	if err != nil {
		account.SentryFileHash = hash
	}

	runWork(account)

	// second run but now with code
	time.Sleep(2 * time.Second)
	account.ShouldRememberPassword = true
	runWork(account)

}

func runWork(account *steam.LogOnDetails) {
	fmt.Println(account, connects)
	client := steam.NewClient()
	steam.InitializeSteamDirectory()
	client.Connect()

	logger := logrus.NewEntry(logrus.New())
	dota := dota2.New(client, logger)
	t := trade.New(client.Web.SessionId, client.Web.SteamLogin, client.Web.SteamLoginSecure, client.SteamId())
	for event := range client.Events() {
		// event,
		fmt.Println(event, reflect.TypeOf(event))
		switch e := event.(type) {
		//If stuck, then restart bot(i think default case is enough???)
		// Ends bot work without disconnect
		case *events.GCConnectionStatusChanged:
			log.Println("Success! State Changed")
			isReady := e.NewState == gcmm.GCConnectionStatus_GCConnectionStatus_HAVE_SESSION
			t.SetReady(isReady)
			if !isReady {
				dota.SayHello()
			}
			return
		// TODO: Read from where LoginKey comes
		// Changed
		case *steam.ConnectedEvent:
			log.Println("Connected! Trying to log in...")
			// Only for testing
			if connects == 1 { // LogOn without email/Steam Guard Code
				client.Auth.LogOn(account)
				connects++
			} else if connects == 2 { // Log on with LoginKey and Hash
				client.Auth.LogOn(account)
			}

			// Release part
			// client.Auth.LogOn(account)

		// Changed
		case *steam.MachineAuthUpdateEvent: // Seems useless (No, it returns hash code for loggin in without email/sean guard codes)
			// check if it works or neded in type converting
			account.SentryFileHash = e.Hash
		// Changed
		case *steam.LoggedOnEvent:
			log.Println("Logged in!")
			fmt.Println(account)
			client.Social.SetPersonaState(steamlang.EPersonaState_Online)
			dota.SetPlaying(true)
			dota.SayHello()
			time.Sleep(1 * time.Second)
		// Deliting password to use the LoginKey
		case *steam.LoginKeyEvent:
			account.Password = ""
			fmt.Println(e.LoginKey)
			account.LoginKey = e.LoginKey
			// Put UniqueID in logs if you want
			_ = e.UniqueId
		case *steam.LogOnFailedEvent:
			if e.Result.String() == "EResult_AccountLogonDeniedVerifiedEmailRequired" || e.Result.String() == "EResult_AccountLogonDenied" {
				var emailCode string
				fmt.Println("Enter the email code")
				fmt.Scanf("%s", &emailCode) // Receiveng users code from email
				fmt.Println(emailCode)
				account.AuthCode = emailCode
			} else if e.Result.String() == "EResult_AccountLogonDeniedNeedTwoFactorCode" {
				var steamGuardCode string
				fmt.Println("Enter the Steam Guard code")
				fmt.Scanf("%s", &steamGuardCode) // Receiveng users code from Steam Guard
				fmt.Println(steamGuardCode)
				account.TwoFactorCode = steamGuardCode
			}
		case *steam.ClientCMListEvent:
			log.Println("Updating Steam IP list!")
			steamIPs := []string{}
			for i := 0; i < len(e.Addresses); i++ {
				steamIPs = append(steamIPs, e.Addresses[i].String())
			}
			os.WriteFile("steam/steamservers.conf", []byte(strings.Join(steamIPs, "\n")), 0666)
		// changed
		case *steam.DisconnectedEvent:
			log.Println("Disconnected!")
			return
		case error:
			log.Printf("Oops! Found an error : %s", e)
		default:
			dota.SayHello()
		}
	}
}

// caseloop:
// 	for {
// 		f.Println("\nEnter case:")
// 		f.Println("[0] - Disconnect")
// 		f.Println("[1] - My STEAMID")
// 		f.Println("[2] - Create new lobby")
// 		f.Println("[3] - Destroy current lobby")
// 		f.Println("[4] - Leave lobby")
// 		f.Println("[5] - Flip teams")
// 		f.Println("[6] - Launch")
// 		f.Println("[7] - Request match - 6897381333")
// 		f.Println("[8] - Invite user")
// 		var i int
// 		f.Scan(&i)
// 		switch i {
// 		case 0:
// 			disconnecting(client)
// 			break caseloop
// 		case 1:
// 			f.Println("Your STEAMID:", client.SteamId())
// 		case 2:
// 			f.Print("Creating Lobby... \n")
// 			lobbySettings := new(gcmm.CMsgPracticeLobbySetDetails)
// 			lobbyData(lobbySettings)
// 			err := dota.LeaveCreateLobby(context.Background(), lobbySettings, true)
// 			if err != nil {
// 				f.Printf("Oops! Found an error : %s", err)
// 				break
// 			}
// 			f.Print("Success!")
// 			id, steamiderror := steamid.NewId("76561198046743792")
// 			if steamiderror != nil {
// 				f.Println(steamiderror)
// 				f.Println("---")
// 			} else {
// 				f.Print("Success!")
// 			}
// 			dota.InviteLobbyMember(id)
// 			f.Println(lobbySettings.GetAllowSpectating())

// 			ctx := context.Background()
// 			response, err := dota.LeaveTeam(ctx, 0)
// 			if err != nil {
// 				f.Println(err)
// 				f.Println("---")
// 			} else {
// 				f.Print(response)
// 			}
// 		// case 2:
// 		// 	f.Print("Creating Lobby... \n")
// 		// 	lobbySettings := new(gcmm.CMsgPracticeLobbySetDetails)

// 		// 	dota.CreateLobby(lobbySettings)
// 		// 	dota.InviteLobbyMember(id)

// 		// 	ctx := context.Background()
// 		// 	team := new(gcmm.CMsgDOTACreateTeam)
// 		// 	name := "team1"
// 		// 	team.Name = &name
// 		// 	response, err := dota.CreateTeam(ctx, team)
// 		// 	if err != nil {
// 		// 		f.Println(err)
// 		// 		f.Println("---")
// 		// 		f.Println(response)
// 		// 		f.Println("---")
// 		// 	}
// 		case 3:
// 			ctx := context.Background()
// 			response, err := dota.DestroyLobby(ctx)
// 			if err != nil {
// 				f.Println(err)
// 				f.Println("---")
// 				f.Println(response)
// 				f.Println("---")
// 			}
// 		case 4:
// 			dota.LeaveLobby()
// 		case 5:
// 			//dota.InvitePrivateChatMember("MyClanChat", 445503465)
// 			dota.FlipLobbyTeams()
// 		case 6:
// 			dota.LaunchLobby()
// 		case 7:
// 			ctx := context.Background()
// 			response, err := dota.RequestMatchDetails(ctx, 6897381333)
// 			if err != nil {
// 				f.Println(err)
// 				f.Println("---")
// 				f.Println(response)
// 				f.Println("---")
// 			} else {
// 				f.Println(response)
// 			}
// 		case 8:
// 			addUser(dota)
// 		case 9:
// 			log.Println("\nWrong key!")
// 		default:
// 			log.Println("\nWrong key!")
// 			//log.Print(logger)
// 		}
// 	}
// }

func info() {
	f.Println("Dota 2 bot")
	f.Println("https://github.com/Saph1s")
	f.Print("Version: ", botVersion, botBuild, "\n\n")
}

func disconnecting(client *steam.Client) {
	client.Disconnect()
	f.Println("Disconnected!")
}

func lobbyData(lobbySettings *gcmm.CMsgPracticeLobbySetDetails) {
	var (
		i         int
		passkey   string
		name      string
		region    uint32
		mode      uint32
		id        uint64
		allowSpec bool
	)
	f.Print("Passkey:")
	f.Scan(&passkey)
	f.Print("Lobby name:")
	f.Scan(&name)
	f.Print("Region (see docs):")
	f.Scan(&region)
	f.Print("Game mode:")
	f.Scan(&mode)
	f.Print("Lobby ID:")
	f.Scan(&id)
	f.Print("Invisible? | [1] - Yes, [2 or other] - No:")
	f.Scan(&i)
	allowSpec = true
	if i == 1 {
		lobbySettings.Visibility = gcmm.DOTALobbyVisibility_DOTALobbyVisibility_Unlisted.Enum()
	} else {
		lobbySettings.Visibility = gcmm.DOTALobbyVisibility_DOTALobbyVisibility_Public.Enum()
	}
	lobbySettings.PassKey = &passkey
	lobbySettings.GameName = &name
	lobbySettings.ServerRegion = &region
	lobbySettings.GameMode = &mode
	lobbySettings.LobbyId = &id
	lobbySettings.AllowSpectating = &allowSpec

	f.Print("Want to add teams? | [1] - Yes, [2 or other] - No:")
	f.Scan(&i)
	if i == 1 {
		f.Println("TODO")
	}

}

// func addUser(dota *dota2.Dota2) {
// 	var cin string
// 	f.Println("Print player STEAMID64")
// 	f.Scan(&cin)
// 	id, steamiderror := steamid.NewId(cin)
// 	if steamiderror != nil {
// 		f.Println(steamiderror)
// 		f.Println("---")
// 	} else {
// 		dota.InviteLobbyMember(id)
// 	}

// }

func readJSON() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		f.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	for key, value := range result {
		f.Println(key, value)
		f.Println("---")
	}
}
