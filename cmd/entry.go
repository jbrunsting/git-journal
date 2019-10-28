package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	MAX_TITLE_LEN  = 80
	DEFAULT_EDITOR = "vim"
)

// entryCmd represents the entry command
var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Create a new journal entry",
	Long:  "Add a new journal entry to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		journalPath := path.Clean(viper.Get("journal-dir").(string))
		journalName := viper.Get("name").(string) + ".md"
		journalFilePath := fmt.Sprintf("%v/%v", journalPath, journalName)

		tempFile, err := ioutil.TempFile(os.TempDir(), "*")
		if err != nil {
			log.Fatal(err)
		}
		tempFilePath := tempFile.Name()
		defer os.Remove(tempFilePath)

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = DEFAULT_EDITOR
		}

		executable, err := exec.LookPath(editor)
		if err != nil {
			log.Fatal(err)
		}

		editorCmd := exec.Command(executable, tempFilePath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		err = editorCmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadFile(tempFilePath)
		if err != nil {
			log.Fatal(err)
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		title := strings.Split(strings.Split(string(body), "\n")[0], ".")[0]
		if len(title) >= MAX_TITLE_LEN {
			title = title[:MAX_TITLE_LEN]
		}
		entry := fmt.Sprintf("## %v ##\n\n%v\n\n", timestamp, string(body))

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
