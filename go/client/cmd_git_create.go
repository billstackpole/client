package client

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
)

type CmdGitCreate struct {
	libkb.Contextified
	repoName keybase1.GitRepoName
	teamName keybase1.TeamName
}

func newCmdGitCreate(cl *libcmdline.CommandLine, g *libkb.GlobalContext) cli.Command {
	return cli.Command{
		Name:         "create",
		ArgumentHelp: "<repo name> [--team=<team name>]",
		Usage:        "Create a personal or team git repository.",
		Description:  "`keybase git create reponame` will create a personal git repo.\n`keybase git create reponame --team=treehouse` will\ncreate a team git repo for the `treehouse` team.\n",
		Action: func(c *cli.Context) {
			cmd := NewCmdGitCreateRunner(g)
			cl.ChooseCommand(cmd, "create", c)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "team",
				Usage: "keybase team name (optional)",
			},
		},
	}
}

func NewCmdGitCreateRunner(g *libkb.GlobalContext) *CmdGitCreate {
	return &CmdGitCreate{Contextified: libkb.NewContextified(g)}
}

func (c *CmdGitCreate) ParseArgv(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		return errors.New("repo name argument required")
	}
	c.repoName = keybase1.GitRepoName(ctx.Args()[0])
	if len(ctx.String("team")) > 0 {
		teamName, err := keybase1.TeamNameFromString(ctx.String("team"))
		if err != nil {
			return err
		}
		c.teamName = teamName
	}

	return nil
}

func (c *CmdGitCreate) Run() error {
	cli, err := GetGitClient(c.G())
	if err != nil {
		return err
	}

	if len(c.teamName.String()) > 0 {
		err = c.runTeam(cli)
	} else {
		err = c.runPersonal(cli)
	}

	if err != nil {
		return err
	}

	dui := c.G().UI.GetDumbOutputUI()
	dui.Printf("Repo created!\n")
	return nil
}

func (c *CmdGitCreate) runPersonal(cli keybase1.GitClient) error {
	if _, err := cli.CreatePersonalRepo(context.Background(), c.repoName); err != nil {
		return err
	}
	return nil
}

func (c *CmdGitCreate) runTeam(cli keybase1.GitClient) error {
	arg := keybase1.CreateTeamRepoArg{
		TeamName: c.teamName,
		RepoName: c.repoName,
	}
	if _, err := cli.CreateTeamRepo(context.Background(), arg); err != nil {
		return err
	}
	return nil
}

func (c *CmdGitCreate) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:    true,
		API:       true,
		KbKeyring: true,
	}
}
