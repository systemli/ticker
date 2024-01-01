package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/storage"

	pwd "github.com/sethvargo/go-password/password"
)

var (
	email        string
	password     string
	isSuperAdmin bool

	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Long:  "Commands for managing users.",
		Args:  cobra.ExactArgs(1),
	}

	userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if email == "" {
				log.Fatal("email is required")
			}
			if password == "" {
				password, err = pwd.Generate(24, 3, 3, false, false)
				if err != nil {
					log.WithError(err).Fatal("could not generate password")
				}
			}

			user, err := storage.NewUser(email, password)
			if err != nil {
				log.WithError(err).Fatal("could not create user")
			}
			user.IsSuperAdmin = isSuperAdmin

			if err := store.SaveUser(&user); err != nil {
				log.WithError(err).Fatal("could not save user")
			}

			fmt.Printf("Created user %d\n", user.ID)
			fmt.Printf("Password: %s\n", password)
		},
	}

	userDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Run: func(cmd *cobra.Command, args []string) {
			if email == "" {
				log.Fatal("email is required")
			}

			user, err := store.FindUserByEmail(email)
			if err != nil {
				log.WithError(err).Fatal("could not find user")
			}

			if err := store.DeleteUser(user); err != nil {
				log.WithError(err).Fatal("could not delete user")
			}

			fmt.Printf("Deleted user %s\n", email)
		},
	}
)

func init() {
	userCmd.AddCommand(userCreateCmd)
	userCreateCmd.Flags().StringVar(&email, "email", "", "email address of the user")
	userCreateCmd.Flags().StringVar(&password, "password", "", "password of the user")
	userCreateCmd.Flags().BoolVar(&isSuperAdmin, "super-admin", false, "make the user a super admin")

	userCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.Flags().StringVar(&email, "email", "", "email address of the user")
}
