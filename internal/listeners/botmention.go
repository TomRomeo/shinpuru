package listeners

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
)

type ListenerBotMention struct {
	config *config.Config

	idLen int32
}

func NewListenerBotMention(container di.Container) *ListenerBotMention {
	return &ListenerBotMention{
		config: container.Get(static.DiConfig).(*config.Config),
		idLen:  0,
	}
}

func (l *ListenerBotMention) Listener(s *discordgo.Session, e *discordgo.MessageCreate) {
	if atomic.LoadInt32(&l.idLen) == 0 {
		atomic.StoreInt32(&l.idLen, int32(len(s.State.User.ID)))
	}

	cLen := int32(len(e.Message.Content))
	if cLen < 3+l.idLen ||
		cLen > 5+l.idLen ||
		e.Message.Content[0] != '<' ||
		e.Message.Content[1] != '@' ||
		e.Author.ID == s.State.User.ID {
		return
	}

	cursor := 2
	if e.Message.Content[2] == '!' {
		cursor = 3
	}

	id := e.Message.Content[cursor : int32(cursor)+l.idLen]
	if id != s.State.User.ID {
		return
	}

	prefix := l.config.Discord.GeneralPrefix
	emb := embedbuilder.New().
		WithColor(static.ColorEmbedDefault).
		WithThumbnail(s.State.User.AvatarURL("64x64"), "", 64, 64).
		WithDescription(fmt.Sprintf("shinpuru Discord Bot v.%s (%s)", embedded.AppVersion, embedded.AppCommit[:6])).
		WithFooter(fmt.Sprintf("© %d Ringo Hoffmann (zekro Development)", time.Now().Year()), "", "").
		AddField("Help", fmt.Sprintf(
			"Type `%shelp` in the chat to get a list of available commands.\n"+
				"You can also use `%shelp <commandInvoke>` to get more details about a command.\n"+
				"[**Here**](https://github.com/zekroTJA/shinpuru/wiki/commands) you can find "+
				"the wiki page with a detailed list of available commands.", prefix, prefix))

	if l.config.WebServer != nil && l.config.WebServer.Enabled {
		emb.AddField("Web Interface", fmt.Sprintf(
			"[**Here**](%s) you can access the web interface.\n"+
				"You can also use the `%slogin` command if you don't want to log in to the web interface via Discord.",
			l.config.WebServer.PublicAddr, prefix))
	}

	emb.AddField("Repository", fmt.Sprintf(
		"[**Here**](https://github.com/zekroTJA/shinpuru) you can find the open source "+
			"repository of shinpuru. Feel free to contribute issues and pull requests, if you want.\n"+
			"You can also use the `%sinfo` command to get more information.", prefix))

	util.SendEmbedRaw(s, e.ChannelID, emb.Build())
}