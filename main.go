package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load("secrets.env")
}

var BotID string
var discord *discordgo.Session
var char *CharacterResponse
var vs *discordgo.VoiceStateUpdate
var BossData *[402]Boss

type Boss struct {
	ID             string
	Name           string
	Image          string
	KilledBossesY  string
	KilledPlayersY string
	KilledBosses   string
	KilledPlayers  string
	LastSeen       string
	Introduced     string
	KilledDaysAgo  int32
	IntervalKill   float64
}

type BossB struct {
	Name     string
	Respawns string
	Type     string
}

type CharacterResponse struct {
	Characters struct {
		Data struct {
			Name              string `json:"name"`
			Title             string `json:"title"`
			Sex               string `json:"sex"`
			Vocation          string `json:"vocation"`
			Level             int    `json:"level"`
			AchievementPoints int    `json:"achievement_points"`
			World             string `json:"world"`
			Residence         string `json:"residence"`
			LastLogin         []struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"last_login"`
			AccountStatus string `json:"account_status"`
			Status        string `json:"status"`
		} `json:"data"`
		Achievements []interface{} `json:"achievements"`
		Deaths       []struct {
			Date struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
			Level    int           `json:"level"`
			Reason   string        `json:"reason"`
			Involved []interface{} `json:"involved"`
		} `json:"deaths"`
		AccountInformation struct {
			LoyaltyTitle string `json:"loyalty_title"`
			Created      struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"created"`
		} `json:"account_information"`
		OtherCharacters []struct {
			Name   string `json:"name"`
			World  string `json:"world"`
			Status string `json:"status"`
		} `json:"other_characters"`
	} `json:"characters"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := discord.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID
	discord.AddHandler(CharacterInfoHandler)
	discord.AddHandler(GetBossesInfoHandler)
	err = discord.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	<-make(chan struct{})
	return
}

func CharacterInfoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	content := strings.ToLower(m.Content)
	if content == "!char" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Musisz podac nazwe postaci np. !char Uther Morlenfra")
	} else if strings.Contains(content, "!char") && len(content) > len("!char") {
		GetCharacterInfo(After(content, "!char"))
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintCharacterInfo())
	}
}

func GetBossesInfoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	content := strings.ToLower(m.Content)
	switch content {
	case "!boss":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Musisz wybrać typ bossów z listy: -poi, -weak, -arch, -rosh, -profit np. !boss -poi")
	case "!boss -poi":
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintBosses("poi"))
	case "!boss -weak":
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintBosses("weak"))
	case "!boss -arch":
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintBosses("arch"))
	case "!boss -rosh":
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintBosses("rosh"))
	case "!boss -profit":
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintBosses("profit"))
	default:
	}
}

// func OnlineWelcomeMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

// 	if m.Author.ID == BotID {
// 		return
// 	}

// 	var channelId string
// 	channels, _ := s.GuildChannels(m.GuildID)
// 	for _, channel := range channels {
// 		if channel.Name == "Tibia" {
// 			channelId = channel.ID
// 		}
// 	}
// 	ch, _ := s.Channel(channelId)
// 	onlineUsers := ch.Recipients
// 	//fmt.Println(vs.UserID)
// 	//fmt.Println(channelId)

// 	_, _ = s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(onlineUsers)))

// }

func GetCharacterInfo(CharName string) {
	response, err := http.Get("https://api.tibiadata.com/v2/characters/" + CharName + ".json")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal([]byte(body), &char)
}

func PrintCharacterInfo() string {

	BasicCharInfo :=
		"```apache" +
			"\nName: " + char.Characters.Data.Name +
			" " + strings.ToUpper(char.Characters.Data.Status) +
			"\nTitle:" + char.Characters.Data.Title +
			"\nSex:" + char.Characters.Data.Sex +
			"\nVocation:" + char.Characters.Data.Vocation +
			"\nLevel:" + strconv.Itoa(char.Characters.Data.Level) +
			"\nAchievement Points:" + strconv.Itoa(char.Characters.Data.AchievementPoints) +
			"\nWorld:" + char.Characters.Data.World +
			"\nResidence:" + char.Characters.Data.Residence +
			"\nAccount Status:" + char.Characters.Data.AccountStatus +
			"```"
	if char.Characters.Data.Name == "" {
		return "Character does not exist."
	}

	if len(char.Characters.Deaths) != 0 {
		CharacterDeaths := "\n"

		for _, death := range char.Characters.Deaths {
			CharacterDeaths += string(death.Date.Date) + ", " + string(death.Date.Timezone) + " " + death.Reason + " at Level " + strconv.Itoa(death.Level) + ".\n"
		}

		CharacterDeaths = "```cs" + "\n" + CharacterDeaths + "\n```"

		char = nil
		return BasicCharInfo + CharacterDeaths
	}

	char = nil
	return BasicCharInfo
}

func After(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}

	adjustedPos := pos + len(a)

	if adjustedPos >= len(value) {
		return ""
	}

	return value[adjustedPos:len(value)]
}

func ScrapWebsite() {
	doc, err := htmlquery.LoadURL("https://guildstats.eu/bosses?monsterName=&world=Monza&rook=0")
	if err != nil {
		fmt.Println(err.Error())
	}
	list := htmlquery.Find(doc, "//td")
	var bN = 0
	var BossData [402]Boss
	for i := 2; i < len(list); i = i + 9 {
		BossData[bN].ID = htmlquery.InnerText(list[i])
		BossData[bN].Name = htmlquery.InnerText(list[i+1])
		BossData[bN].Image = htmlquery.InnerText(list[i+2])
		BossData[bN].KilledBossesY = htmlquery.InnerText(list[i+3])
		BossData[bN].KilledPlayersY = htmlquery.InnerText(list[i+4])
		BossData[bN].KilledBosses = htmlquery.InnerText(list[i+5])
		BossData[bN].KilledPlayers = htmlquery.InnerText(list[i+6])
		BossData[bN].LastSeen = htmlquery.InnerText(list[i+7])
		BossData[bN].Introduced = htmlquery.InnerText(list[i+8])

		lastDateWhenBossKilled, err := time.Parse("2006-01-02", BossData[bN].LastSeen)
		if err != nil {
			fmt.Println(err.Error())
		}
		daysFromKill := time.Now().Sub(lastDateWhenBossKilled).Hours() / 24
		BossData[bN].KilledDaysAgo = int32(daysFromKill)

		introducedDate, err := time.Parse("2006-01-02", BossData[bN].Introduced)
		if err != nil {
			fmt.Println(err.Error())
		}
		daysFromIntroduced := time.Now().Sub(introducedDate).Hours() / 24

		pr, err := strconv.ParseFloat(BossData[bN].KilledBosses, 64)
		if err != nil {
			fmt.Println(err.Error())
		}

		if pr == 0 {
			BossData[bN].IntervalKill = 0
			//fmt.Println(BossData[bN].IntervalKill)
		} else {
			intervalKill := daysFromIntroduced / pr
			BossData[bN].IntervalKill = intervalKill
			//fmt.Println(BossData[bN].IntervalKill)
		}
		var x float64 = float64(BossData[bN].KilledDaysAgo)

		c := x - BossData[bN].IntervalKill

		fmt.Print(BossData[bN].Name)
		fmt.Print(" Boss killed: ", BossData[bN].KilledDaysAgo)
		fmt.Print(" Interval kill: ", BossData[bN].IntervalKill)
		fmt.Println(" Prob: ", c)

		bN++
	}

	fmt.Println(BossData[401].Name + " Last Seen: " + BossData[401].LastSeen)

}

func PrintBosses(bossType string) string {

	var PoiBoss, WeakBoss, RoshamuulBoss, ArchdemonsBoss, ProfitBoss = scrapTibiaBosses("monza")
	PoiBossInfo := "```css\n"
	for _, Pb := range PoiBoss {
		PoiBossInfo += Pb.Name + "," + Pb.Respawns + "\n"
	}
	PoiBossInfo += "```"
	PoiBossInfo = strings.Replace(PoiBossInfo, "Low Chance", "[Low Chance]", -1)
	PoiBossInfo = strings.Replace(PoiBossInfo, "No Chance", "[No Chance]", -1)

	WeakBossInfo := "```css\n"
	for _, Wb := range WeakBoss {
		WeakBossInfo += Wb.Name + "," + Wb.Respawns + "\n"
	}
	WeakBossInfo += "```"
	WeakBossInfo = strings.Replace(WeakBossInfo, "Low Chance", "[Low Chance]", -1)
	WeakBossInfo = strings.Replace(WeakBossInfo, "No Chance", "[No Chance]", -1)

	RoshamuulBossInfo := "```css\n"
	for _, Rb := range RoshamuulBoss {
		RoshamuulBossInfo += Rb.Name + "," + Rb.Respawns + "\n"
	}
	RoshamuulBossInfo += "```"
	RoshamuulBossInfo = strings.Replace(RoshamuulBossInfo, "Low Chance", "[Low Chance]", -1)
	RoshamuulBossInfo = strings.Replace(RoshamuulBossInfo, "No Chance", "[No Chance]", -1)

	ArchdemonsBossInfo := "```css\n"
	for _, Ab := range ArchdemonsBoss {
		ArchdemonsBossInfo += Ab.Name + "," + Ab.Respawns + "\n"
	}
	ArchdemonsBossInfo += "```"
	ArchdemonsBossInfo = strings.Replace(ArchdemonsBossInfo, "Low Chance", "[Low Chance]", -1)
	ArchdemonsBossInfo = strings.Replace(ArchdemonsBossInfo, "No Chance", "[No Chance]", -1)

	ProfitBossInfo := "```css\n"
	for _, Pb := range ProfitBoss {
		ProfitBossInfo += Pb.Name + "," + Pb.Respawns + "\n"
	}
	ProfitBossInfo += "```"
	ProfitBossInfo = strings.Replace(ProfitBossInfo, "Low Chance", "[Low Chance]", -1)
	ProfitBossInfo = strings.Replace(ProfitBossInfo, "No Chance", "[No Chance]", -1)

	switch bossType {
	case "poi":
		return PoiBossInfo
	case "weak":
		return WeakBossInfo
	case "rosh":
		return RoshamuulBossInfo
	case "arch":
		return ArchdemonsBossInfo
	case "profit":
		return ProfitBossInfo
	default:
		return "Error!"
	}
}

func scrapTibiaBosses(world string) ([7]BossB, [15]BossB, [2]BossB, [4]BossB, [11]BossB) {
	doc, err := htmlquery.LoadURL("https://www.tibiabosses.com/" + world + "/")
	if err != nil {
		fmt.Println(err.Error())
	}
	var PoiBoss [7]BossB
	var WeakBoss [15]BossB
	var RoshamuulBoss [2]BossB
	var ArchdemonsBoss [4]BossB
	var ProfitBoss [11]BossB
	PoiBNames := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[5]/div[2]/div/div/a/@href`)
	PoiB := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[5]/div[2]/div/div`)

	PoiBLines := strings.Split(htmlquery.InnerText(PoiB[0]), "\n")

	for i := 0; i < len(PoiBoss); i++ {
		PoiBoss[i].Name = strings.Split(htmlquery.InnerText(PoiBNames[i]), "/bossopedia/")[1]
		PoiBoss[i].Respawns = PoiBLines[i+1]
		PoiBoss[i].Type = "POI Bosses"
	}

	ProfitBNames := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[3]/div[2]/div/div/a/@href`)
	ProfitB := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[3]/div[2]/div/div`)

	ProfitBLines := strings.Split(htmlquery.InnerText(ProfitB[0]), "   ")

	for i := 0; i < len(ProfitBoss); i++ {
		ProfitBoss[i].Name = strings.Split(htmlquery.InnerText(ProfitBNames[i]), "/bossopedia/")[1]
		ProfitBoss[i].Type = "Profit Boss"
	}

	ProfitBoss[0].Respawns = ProfitBLines[0]
	ProfitBoss[1].Respawns = ProfitBLines[1]
	ProfitBoss[2].Respawns = ProfitBLines[2]
	ProfitBoss[3].Respawns = strings.Replace((ProfitBLines[3] + " " + ProfitBLines[5] + ProfitBLines[6] + " " + ProfitBLines[8]), "\n", "", -1)
	ProfitBoss[4].Respawns = ProfitBLines[11] + " " + ProfitBLines[12]
	ProfitBoss[5].Respawns = ProfitBLines[15] + " " + ProfitBLines[16]
	ProfitBoss[6].Respawns = ProfitBLines[19] + " " + ProfitBLines[20]
	ProfitBoss[7].Respawns = ProfitBLines[23] + " " + ProfitBLines[24]
	ProfitBoss[8].Respawns = ProfitBLines[27] + " " + ProfitBLines[28]
	ProfitBoss[9].Respawns = ProfitBLines[31] + " " + ProfitBLines[32]
	ProfitBoss[10].Respawns = ProfitBLines[35] + " " + ProfitBLines[36]

	ArchdemonsBNames := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[6]/div[2]/div/div/a/@href`)
	ArchdemonsB := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[6]/div[2]/div/div`)
	ArchdemonsBLines := strings.Split(htmlquery.InnerText(ArchdemonsB[0]), "\n")

	for i := 0; i < len(ArchdemonsBoss); i++ {
		ArchdemonsBoss[i].Name = strings.Split(htmlquery.InnerText(ArchdemonsBNames[i]), "/bossopedia/")[1]
		ArchdemonsBoss[i].Respawns = ArchdemonsBLines[i+1]
		ArchdemonsBoss[i].Type = "Archdemon"
	}

	RoshamuulBNames := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[6]/div[1]/div/div/a/@href`)
	RoshamuulB := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[6]/div[1]/div/div`)
	RoshamuulBLines := strings.Split(htmlquery.InnerText(RoshamuulB[0]), "\n")
	for i := 0; i < len(RoshamuulBoss); i++ {
		RoshamuulBoss[i].Name = strings.Split(htmlquery.InnerText(RoshamuulBNames[i]), "/bossopedia/")[1]
		RoshamuulBoss[i].Respawns = RoshamuulBLines[i+1]
		RoshamuulBoss[i].Type = "Roshamuul Boss"
	}

	WeakBNames := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[5]/div[1]/div/div/a/@href`)
	WeakB := htmlquery.Find(doc, `/html/body/div[1]/div[1]/section/article/div/div/div[1]/div[5]/div[1]/div/div`)
	WeakBLines := strings.Split(htmlquery.InnerText(WeakB[0]), "\n")
	for i := 0; i < len(WeakBoss); i++ {
		WeakBoss[i].Name = strings.Split(htmlquery.InnerText(WeakBNames[i]), "/bossopedia/")[1]
		WeakBoss[i].Respawns = WeakBLines[i+1]
		WeakBoss[i].Type = "Weak Boss"
	}

	return PoiBoss, WeakBoss, RoshamuulBoss, ArchdemonsBoss, ProfitBoss
}
