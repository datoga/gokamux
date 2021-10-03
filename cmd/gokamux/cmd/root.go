package cmd

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/datoga/saramaconfig"
)

var cfgFile string
var Verbose bool
var Version = "1.0.0" //TODO: Get from arg in build

type Config struct {
	Brokers []string
	Sarama  sarama.Config
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gokamux",
	Version: Version,
	Short:   "gokamux is a service to mux Kafka streams",
	Long: `gokamux takes different input Kafka streams and can perform different actions on them:
	- Filtering (pass or restrict messages according different criteria).
	- Transforming (read an entry, executes some transformation and passes the result).
	- Output (writes the result in a different output topics).
You can define declaratively the Kafka config and your setup, or provide your own filters and transformers via the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed loading the config with error", err)
			return
		}

		if Verbose {
			fmt.Printf("Config: %+v\n", *cfg)
		}
	},
}

func loadConfig() (*Config, error) {
	brokers := viper.GetStringSlice("brokers")

	if len(brokers) == 0 {
		brokers = []string{"localhost:9092"}
	}

	saramaCfg, err := saramaconfig.NewFromViper(viper.GetViper())

	if err != nil {
		return nil, fmt.Errorf("failed configuring Sarama with error %v", err)
	}

	return &Config{
		Brokers: brokers,
		Sarama:  *saramaCfg,
	}, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "c", "config file (default is $HOME/.gokamux.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gokamux" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/gokamux/")
		viper.SetConfigType("toml")
		viper.SetConfigName("gokamux")
	}

	viper.SetEnvPrefix("gokamux")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if Verbose {
		fmt.Printf("Config file selected: %s\n", viper.ConfigFileUsed())
	}
}
