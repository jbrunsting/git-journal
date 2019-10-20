package cmd

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"log"
    "path"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the journal",
	Long:  "Create the repository to store the journal entries",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
        journalPath := path.Clean(viper.Get("journal-dir").(string))
		fmt.Printf("Creating journal directory %v\n", journalPath)
		err = os.Mkdir(journalPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

        repo, err := git.PlainInit(journalPath + "/.git", true)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("%v\n", repo)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
