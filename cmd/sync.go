package cmd

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"log"
	"path"

	"github.com/spf13/cobra"
)

var remoteUrl string
var keyFile string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with a remote",
	Long:  "Pull changes from a remote repository and submit updates",
	Run: func(cmd *cobra.Command, args []string) {
		pem, err := ioutil.ReadFile(keyFile)
		if err != nil {
			log.Fatal(err)
		}

		signer, err := ssh.ParsePrivateKey(pem)
		if err != nil {
			log.Fatal(err)
		}

		auth := &gitssh.PublicKeys{User: "git", Signer: signer}

		journalPath := path.Clean(viper.Get("journal-dir").(string))
		r, err := git.PlainOpen(journalPath)
		if err != nil {
			log.Fatal(err)
		}

		err = r.DeleteRemote("origin")
		if err != nil {
			log.Fatal(err)
		}

		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{remoteUrl},
		})
		if err != nil {
			log.Fatal(err)
		}

		w, err := r.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       auth,
		})
		if err != nil && err.Error() != "remote repository is empty" {
			log.Fatal(err)
		}

		err = r.Push(&git.PushOptions{Auth: auth})
        if err != nil {
            log.Fatal(err)
        }
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.PersistentFlags().StringVar(&remoteUrl, "remote-url", "", "Remote repository URL")
	syncCmd.PersistentFlags().StringVar(&keyFile, "key-file", "", "File containing the SSH key")
}
