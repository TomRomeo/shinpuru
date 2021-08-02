package commands

import (
	"fmt"
	"strconv"

	"github.com/zekroTJA/shinpuru/pkg/dmdialog"
	"github.com/zekroTJA/shireikan"
)

type CmdTest struct {
}

func (c *CmdTest) GetInvokes() []string {
	return []string{"test"}
}

func (c *CmdTest) GetDescription() string {
	return "Just for testing purposes."
}

func (c *CmdTest) GetHelp() string {
	return ""
}

func (c *CmdTest) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *CmdTest) GetDomainName() string {
	return "sp.test"
}

func (c *CmdTest) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdTest) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdTest) Exec(ctx shireikan.Context) (err error) {
	a, err := dmdialog.New(ctx.GetSession()).
		AddQuestion(dmdialog.Question{
			ID:   "name",
			Text: "Whats your name?",
		}).
		AddQuestion(dmdialog.Question{
			ID:   "age",
			Text: "Whats your age?",
			Validator: func(s string) error {
				_, err := strconv.Atoi(s)
				return err
			},
		}).
		Send(ctx.GetMember().User.ID)

	if err != nil {
		return
	}

	fmt.Println(a.Await())

	return
}
