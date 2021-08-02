package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/dmdialog"
	"github.com/zekroTJA/shireikan"
)

const (
	apiIDLen  = 32
	apiKeyLen = 64
)

type CmdExec struct {
}

func (c *CmdExec) GetInvokes() []string {
	return []string{"exec", "ex", "execute", "jdoodle"}
}

func (c *CmdExec) GetDescription() string {
	return "Setup code execution of code embeds."
}

func (c *CmdExec) GetHelp() string {
	return "`exec setup` - enter jdoodle setup\n" +
		"`exec reset` - disable and delete token from database\n" +
		"`exec toggle` - toggle the enabled state of the code execution\n" +
		"`exec` - display the current code exec status"
}

func (c *CmdExec) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdExec) GetDomainName() string {
	return "sp.chat.exec"
}

func (c *CmdExec) GetSubPermissionRules() []shireikan.SubPermission {
	return []shireikan.SubPermission{
		{
			Term:        "exec",
			Explicit:    false,
			Description: "Allows activating a code execution in chat via reaction",
		},
	}
}

func (c *CmdExec) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdExec) Exec(ctx shireikan.Context) error {

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
	case "toggle":
		return c.setState(ctx)
	case "enable", "on":
		return c.setState(ctx, true)
	case "disable", "off":
		return c.setState(ctx, false)
	case "setup":
		return c.setup(ctx)
	case "reset":
		return c.reset(ctx)
	default:
		return c.status(ctx)
	}

	return nil
}

func (c *CmdExec) setState(ctx shireikan.Context, enable ...bool) (err error) {
	db := ctx.GetObject(static.DiDatabase).(database.Database)

	state, err := db.GetGuildExec(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if state.Provider == "" {
		err = util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Code execution is not set up at the moment. Use `exec setup` to set up code execution before toggling it.").
			DeleteAfter(12 * time.Second).Error()
		return
	}

	if len(enable) == 0 {
		state.Enabled = !state.Enabled
	} else {
		state.Enabled = enable[0]
	}

	if err = db.SetGuildExec(ctx.GetGuild().ID, state); err != nil {
		return
	}

	var emb *util.EmbedMessage
	if state.Enabled {
		emb = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Code execution is now enabled.", "", static.ColorEmbedGreen)
	} else {
		emb = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Code execution is now disabled.", "", static.ColorEmbedOrange)
	}

	return emb.DeleteAfter(6 * time.Second).Error()
}

func (c *CmdExec) setup(ctx shireikan.Context) (err error) {
	db := ctx.GetObject(static.DiDatabase).(database.Database)

	ans, err := dmdialog.New(ctx.GetSession()).
		AddQuestion(dmdialog.Question{
			ID:   "service",
			Text: "Which service do you want to use?\nAvailable: `ranna`, `jdoodle`",
			Validator: func(s string) error {
				s = strings.ToLower(s)
				if s != "ranna" && s != "jdoodle" {
					return errors.New("Invalid service.")
				}
				return nil
			},
			Formatter: strings.ToLower,
		}).
		Send(ctx.GetUser().ID)
	if err != nil {
		return
	}

	err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"The setup is performed via DMs because secrets might be needed to be entered. Take a look into your DMs.", "", 0).
		DeleteAfter(6 * time.Second).Error()
	if err != nil {
		return
	}

	res, m := ans.Await()
	if res != dmdialog.ResultOK {
		return
	}

	provider := m["service"]
	fmt.Println(provider)
	var cfg interface{}

	switch provider {

	case "ranna":
		ans, err = dmdialog.New(ctx.GetSession()).
			AddQuestion(dmdialog.Question{
				ID:   "endpoint",
				Text: "Which ranna endpoint should be used?\nExample: `https://public.ranna.zekro.de`",
				Validator: func(s string) error {
					if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
						return errors.New("The endpoint URL must either start with `http://` or `https://`")
					}
					return nil
				},
			}).
			AddQuestion(dmdialog.Question{
				ID: "token",
				Text: "Do you want to use a token for this endpoint? When yes, please enter the token now. Otherwise, enter `no`.\n" +
					"After entering the token, you can safely delete the message containing the token due to safety measures.",
				Formatter: func(s string) string {
					if strings.ToLower(s) == "no" {
						return ""
					}
					return s
				},
			}).
			Send(ctx.GetUser().ID)
		if err != nil {
			return
		}
		res, m = ans.Await()
		if res != dmdialog.ResultOK {
			return
		}
		cfg = models.ExecConfigRanna{
			Endpoint: m["endpoint"],
			Token:    m["token"],
		}

	case "jdoodle":
		ans, err = dmdialog.New(ctx.GetSession()).
			AddQuestion(dmdialog.Question{
				ID:   "id",
				Text: "First of all, we need your JDoodle `client ID`. You can obtain it from here:\nhttps://www.jdoodle.com/compiler-api/",
			}).
			AddQuestion(dmdialog.Question{
				ID: "secret",
				Text: "Now, we need your `Client Secret` which can also be obtained in the page mentioned above.\n" +
					"After entering the secret, you can safely delete the message containing the secret due to safety measures.",
			}).
			Send(ctx.GetUser().ID)
		if err != nil {
			return
		}
		res, m = ans.Await()
		if res != dmdialog.ResultOK {
			return
		}
		cfg = models.ExecConfigJdoodle{
			ClientID:     m["id"],
			ClientSecret: m["secret"],
		}
	}

	ans, err = dmdialog.New(ctx.GetSession()).
		AddQuestion(dmdialog.Question{
			ID: "enable",
			Text: "Okay, last question: Do you want to enable code execution on the guild? You can always enable/disable it using the `exec toggle` command.\n" +
				"Type `yes` to enable code execution or anything else to leave it disabled.",
			Formatter: strings.ToLower,
		}).
		SetFinishMessage("☑️ Code execution is now set up.").
		Send(ctx.GetUser().ID)
	if err != nil {
		return
	}
	res, m = ans.Await()
	if res != dmdialog.ResultOK {
		return
	}

	state := &models.ExecConfig{
		Provider: provider,
	}
	state.Enabled = m["enable"] == "yes"

	if err = state.SerializeConfig(cfg); err != nil {
		return
	}

	if err = db.SetGuildExec(ctx.GetGuild().ID, state); err != nil {
		return
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Successfully configured code execution.", "", 0).DeleteAfter(6 * time.Second).Error()
}

func (c *CmdExec) reset(ctx shireikan.Context) (err error) {
	db := ctx.GetObject(static.DiDatabase).(database.Database)

	if err = db.SetGuildExec(ctx.GetGuild().ID, &models.ExecConfig{}); err != nil {
		return
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Successfully reset code execution configuration and removed all entered secrets.", "", 0).
		DeleteAfter(6 * time.Second).Error()
}

func (c *CmdExec) status(ctx shireikan.Context) (err error) {
	db := ctx.GetObject(static.DiDatabase).(database.Database)

	state, err := db.GetGuildExec(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if state.Provider == "" {
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Code execution is currently not set up on this guild.\n"+
				"Use `exec setup` to set up code execution.", "", 0).
			DeleteAfter(10 * time.Second).Error()
	}

	msg := fmt.Sprintf("Enabled: `%t`\nProvider: `%s`", state.Enabled, state.Provider)
	clr := static.ColorEmbedOrange
	if state.Enabled {
		clr = static.ColorEmbedGreen
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		msg, "", clr).DeleteAfter(10 * time.Second).Error()
}
