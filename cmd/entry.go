package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	MAX_TITLE_LEN = 80
)

// entryCmd represents the entry command
var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Create a new journal entry",
	Long:  "Add a new journal entry to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("entry called. TODO: Open up default editor and let user write something")

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		body := strings.TrimSpace(`This is a temp message
And other things. Temp message and stuff. And things.

Ando ther things on a new paragraph`)
		title := strings.Split(strings.Split(body, "\n")[0], ".")[0]
        if len(title) >= MAX_TITLE_LEN {
            title = title[:MAX_TITLE_LEN]
        }
		entry := fmt.Sprintf("## %v ##\n\n%v\n\n", timestamp, body)

		journalPath := path.Clean(viper.Get("journal-dir").(string))
		journalName := viper.Get("name").(string) + ".md"
		journalFilePath := fmt.Sprintf("%v/%v", journalPath, journalName)

		f, err := os.OpenFile(journalFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.Write([]byte(entry))
		if err != nil {
			log.Fatal(err)
		}

		f.Close()

		r, err := git.PlainOpen(journalPath)
		if err != nil {
			log.Fatal(err)
		}

		w, err := r.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Add(journalName)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Commit(title, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "TODO",
				Email: "TODO",
				When:  time.Now(),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(entryCmd)
}
