package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		apiFlag, err := cmd.Flags().GetBool("api")
		if err != nil {
			fmt.Println(err)
		}

		moduleFlag, err := cmd.Flags().GetBool("module")
		if err != nil {
			fmt.Println(err)
		}

		if moduleFlag {
			moduleName, err := cmd.Flags().GetString("name")
			if err != nil {
				fmt.Println(err)
			}

			authFlag, err := cmd.Flags().GetBool("auth")
			if err != nil {
				fmt.Println(err)
			}

			if authFlag {
				fmt.Printf("Creating auth module project with name: %s\n", moduleName)
				// Add your logic for creating an auth module project here
			} else {
				fmt.Printf("Creating CRUD module project with name: %s\n", moduleName)
				// Add your logic for creating a CRUD module project here
			}
		}

		if apiFlag {
			hexagonFlag, err := cmd.Flags().GetBool("hexagon")
			if err != nil {
				fmt.Println(err)
			}

			mvcFlag, err := cmd.Flags().GetBool("mvc")
			if err != nil {
				fmt.Println(err)
			}

			if hexagonFlag && mvcFlag {
				fmt.Println("You can only choose one architecture")
				return
			} else if !hexagonFlag && !mvcFlag {
				fmt.Println("Please choose an architecture")
				return
			}

			if hexagonFlag {
				fmt.Println("Creating hexagon architecture project...")
				// Add your logic for creating a hexagon architecture project here
			} else if mvcFlag {
				fmt.Println("Creating MVC architecture project...")
				// Add your logic for creating an MVC architecture project here
			} else {
				fmt.Println("Please specify an architecture to create (hexagon or mvc)")
			}
		}

	},
}

func init() {
	// architecture types
	createCmd.Flags().Bool("api", false, "create API project")
	createCmd.Flags().Bool("hexagon", false, "create hexagon architecture project")
	createCmd.Flags().Bool("mvc", false, "create MVC architecture project")

	// modules types
	createCmd.Flags().Bool("module", false, "create a crud module project")
	createCmd.Flags().String("name", "", "name of the module")
	createCmd.Flags().Bool("auth", false, "create auth module project")

	rootCmd.AddCommand(createCmd)
}
