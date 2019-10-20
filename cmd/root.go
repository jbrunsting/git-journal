package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var journalPath string

var rootCmd = &cobra.Command{
	Use:   "git-journal",
	Short: "A git-backed journal",
	Long:  "A journal which stores entries using git",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("journal-dir", ".journal")
	viper.SetDefault("name", "journal")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-journal.yaml)")
	rootCmd.PersistentFlags().StringVar(&journalPath, "journal-dir", "", "path to the directory with the journal")
	rootCmd.PersistentFlags().StringVar(&journalPath, "name", "", "name of the journal")

	viper.BindPFlag("journal-dir", rootCmd.PersistentFlags().Lookup("journal-dir"))
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".git-journal")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
