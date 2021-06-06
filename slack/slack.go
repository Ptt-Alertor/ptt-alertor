package slack

import (
	"net/http"

	"github.com/meifamily/ptt-alertor/models"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/command"
	"github.com/nlopes/slack"
)

// auth url: https://slack.com/oauth/authorize?client_id=139913746162.341396063444&scope=im:write,bot,commands,chat:write:bot&redirect_uri=https://83c0624e.ngrok.io/slack/auth/redirect

var clientID = "139913746162.341396063444"
var clientSecret = "133bb122166627a235c91bea9b385b22"
var verificationToken = "E69baLlnBPySAz47cLb9Hhfq"
var redirectURI = "https://83c0624e.ngrok.io/slack/auth/redirect"

type Slack struct {
	accessToken string
	channel     string
}

func New(accessToken, channel string) *Slack {
	return &Slack{
		accessToken: accessToken,
		channel:     channel,
	}
}

func HandleAuth(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query().Get("code")
	resp, err := slack.GetOAuthResponse(clientID, clientSecret, code, "", false)
	if err != nil {
		panic(err)
	}
	command.HandleSlackFollow(resp.TeamID, "", resp.AccessToken, resp.Bot.BotAccessToken)

	api := slack.New(resp.AccessToken)
	_, _, channelID, err := api.OpenIMChannel(resp.UserID)
	if err != nil {
		panic(err)
	}
	command.HandleSlackFollow(resp.TeamID, channelID, resp.AccessToken, resp.Bot.BotAccessToken)

	responseText := "歡迎使用 Ptt Alertor\n輸入「指令」查看相關功能。\n\n觀看Demo:\nhttps://media.giphy.com/media/3ohzdF6vidM6I49lQs/giphy.gif"
	s := New(resp.AccessToken, channelID)
	s.SendTextMessage(responseText)
}

func HandleCommand(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cmd, err := slack.SlashCommandParse(r)
	if err != nil {
		panic(err)
	}
	if cmd.Token != verificationToken {
		return
	}

	id := cmd.TeamID + "|" + cmd.ChannelID
	u := models.User.Find(id)
	if u.Profile.Slack.AccessToken == "" {
		t := models.User.Find(cmd.TeamID)
		command.HandleSlackFollow(cmd.TeamID, cmd.ChannelID, t.Profile.Slack.AccessToken, t.Profile.Slack.BotAccessToken)
		u = models.User.Find(id)
	}
	s := New(u.Slack.AccessToken, u.Slack.Channel)
	responseText := command.HandleCommand(cmd.Text, id)
	s.sendHiddenMessage(responseText, cmd.UserID)
}

func (s Slack) sendHiddenMessage(text, userID string) {
	api := slack.New(s.accessToken)
	api.SetDebug(true)
	api.PostEphemeral(s.channel, userID, slack.MsgOptionText(text, false))
}

func (s Slack) SendTextMessage(text string) {
	api := slack.New(s.accessToken)
	api.SetDebug(true)
	api.PostMessage(s.channel, text, slack.PostMessageParameters{
		EscapeText:  false,
		UnfurlLinks: true,
		UnfurlMedia: true,
	})
}
