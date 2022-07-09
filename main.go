package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kevineaton/art/transformer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "ART"
)

func main() {
	rand.Seed(time.Now().Unix())

	root := Root()
	if err := root.Execute(); err != nil {
		fmt.Printf("ERROR: Could not establish the CLI: %+v\n", err)
		os.Exit(1)
	}

}

// Root establishes the root command and sets up the initial hooks for config
func Root() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "art",
		Short: "A small app for trying out generative art",
		Long:  "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// connect viper
			return initializeViper(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ERROR: You must supply a command. Please view the help documentation and try again.\n")
			os.Exit(1)
		},
	}

	rootCmd.AddCommand(transformer.GetCommand())

	return rootCmd
}

func initializeViper(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName("settings")
	v.AddConfigPath(".")

	// attempt to read the file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.SetEnvPrefix(envPrefix)

	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
