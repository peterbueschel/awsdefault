package main

import (
	"fmt"
	"log"
	"os"

	"github.com/peterbueschel/awsdefault"
	"github.com/urfave/cli"
)

func getProfiles(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "profiles", "available"},
		Usage:   "Returns all available profiles from the AWS credentials file.",
		Action: func(c *cli.Context) error {
			names := file.GetProfilesNames()
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}
}

func getUsedProfile(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "get",
		Aliases: []string{"show", "is", "now", "curr"},
		Usage:   "Returns the current AWS default profile.",
		Action: func(c *cli.Context) error {
			n, _, err := file.GetUsedProfileNameAndIndex()
			if err != nil {
				return err
			}
			fmt.Println(n)
			return nil
		},
	}
}

func unsetDefaultProfile(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "unset",
		Aliases: []string{"rm", "stop", "not"},
		Usage:   "unset the AWS default profile.",
		Action: func(c *cli.Context) error {
			if err := file.UnSetDefault(); err != nil {
				return err
			}
			return nil
		},
	}
}

func setDefaultProfile(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "set",
		Aliases: []string{"to", "should", "replace"},
		Usage:   "Set/replace the AWS default profile to a given profile. Requires a profile name.",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf(
					"the name of the profile used to become the new default is required",
				)
			}
			if err := file.SetDefaultTo(c.Args().First()); err != nil {
				return err
			}
			return nil
		},
	}
}

func main() {
	file, err := awsdefault.GetCredentialsFile()
	if err != nil {
		log.Fatalf("[AWSFLIPPER][ERROR] %v.\n", err)
	}

	app := cli.NewApp()

	app.Commands = []cli.Command{
		*setDefaultProfile(file),
		*unsetDefaultProfile(file),
		*getUsedProfile(file),
		*getProfiles(file),
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("[AWSFLIPPER][ERROR] %s.\n", err)
	}
}
