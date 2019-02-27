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
			n, idx, err := file.GetUsedProfileNameAndIndex()
			if err != nil {
				if idx == -2 {
					fmt.Println(n)
					return nil
				}
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
			return file.UnSetDefault()
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
			return file.SetDefaultTo(c.Args().First())
		},
	}
}

func getUsedID(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "id",
		Aliases: []string{"aws_access_key_id"},
		Usage:   "Returns the AWS_ACCESS_KEY_ID of the currently used profile",
		Action: func(c *cli.Context) error {
			id, err := file.GetUsedID()
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}
}

func getUsedKey(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "key",
		Aliases: []string{"aws_secret_access_key"},
		Usage:   "Returns the AWS_SECRET_ACCESS_KEY of the currently used profile",
		Action: func(c *cli.Context) error {
			k, err := file.GetUsedKey()
			if err != nil {
				return err
			}
			fmt.Println(k)
			return nil
		},
	}
}

func printCredential(file *awsdefault.CredentialsFile) *cli.Command {
	return &cli.Command{
		Name:    "export",
		Aliases: []string{"envs"},
		Usage:   "Returns the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY of the currently used profile in form of export commands.",
		Action: func(c *cli.Context) error {
			id, err := file.GetUsedID()
			if err != nil {
				return err
			}
			k, err := file.GetUsedKey()
			if err != nil {
				return err
			}
			fmt.Printf("export AWS_ACCESS_KEY_ID=%s\nexport AWS_SECRET_ACCESS_KEY=%s\n", id, k)
			return nil
		},
	}
}

func main() {
	file, err := awsdefault.GetCredentialsFile()
	if err != nil {
		log.Fatalf("[AWSDEFAULT][ERROR] %v.\n", err)
	}

	app := cli.NewApp()

	app.Commands = []cli.Command{
		*setDefaultProfile(file),
		*unsetDefaultProfile(file),
		*getUsedProfile(file),
		*getUsedID(file),
		*getUsedKey(file),
		*printCredential(file),
		*getProfiles(file),
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("[AWSDEFAULT][ERROR] %s.\n", err)
	}
}
