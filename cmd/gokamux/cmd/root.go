package cmd

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/datoga/gokamux"
	"github.com/datoga/gokamux/modules"
	"github.com/spf13/cobra"

	_ "github.com/datoga/gokamux/modules/builtin/all"

	"github.com/spf13/viper"
)

var cfgFile string
var saramaCfgFile string

var PluginsPath string
var Verbose bool
var Version = "1.0.0" //TODO: Get from arg in build

var Cfg *Config
var SaramaCfg *sarama.Config

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
		muxer := gokamux.NewMuxer().
			Brokers(Cfg.Brokers...).
			Input(Cfg.InputTopics...).
			Output(Cfg.OutputTopics...).
			Step(Cfg.Steps...)

		cobra.CheckErr(muxer.Run(context.Background()))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(configureLog, initConfigGokaMux, initSaramaConfig, discoverPlugins)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "gokamux config file (default is $HOME/.gokamux.toml)")

	rootCmd.PersistentFlags().StringVarP(&saramaCfgFile, "sarama-config", "s", "", "sarama config file (default is $HOME/.sarama.toml) (optional)")

	rootCmd.PersistentFlags().StringVarP(&PluginsPath, "plugins", "p", "plugins", "plugins directory (optional)")

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func configureLog() {
	if !Verbose {
		log.SetOutput(io.Discard)
	}
}

func discoverPlugins() {
	if _, err := os.Stat(PluginsPath); !os.IsNotExist(err) {
		modules.MustDiscoverAndRegister(PluginsPath)
	}

	log.Println(len(modules.List()), "modules loaded")

	for _, m := range modules.List() {
		log.Println("Module", m)
	}
}

// initConfigGokaMux reads in a GokaMux config file and ENV variables if set.
func initConfigGokaMux() {
	v, found := initConfig(cfgFile, "gokamux")

	if !found {
		cobra.CheckErr(errors.New("failed loading the gokamux config"))
	}

	cfg, err := loadGokaMuxConfig(v)

	cobra.CheckErr(err)

	log.Printf("GokaMux Config: %+v\n", *cfg)

	Cfg = cfg
}

// initSaramaConfig reads in Sarama config file and ENV variables if set.
func initSaramaConfig() {
	v, found := initConfig(saramaCfgFile, "sarama")

	if !found {
		log.Println("No Sarama config file found, default config will be used")
		return
	}

	saramaCfg, err := loadSaramaConfig(v)

	cobra.CheckErr(err)

	log.Printf("Sarama Config: %+v\n", *saramaCfg)

	SaramaCfg = saramaCfg
}

// initConfig reads in a config file and ENV variables if set.
func initConfig(cfgFile string, cfgName string) (*viper.Viper, bool) {
	v := viper.New()

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".$CONFIGNAME" (without extension).
		v.AddConfigPath(home)
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/gokamux/")
		v.SetConfigType("toml")
		v.SetConfigName(cfgName)

		v.SetEnvPrefix(cfgName)
		v.AutomaticEnv() // read in environment variables that match
	}

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		log.Printf("Using config file for %s: %s\n", cfgName, v.ConfigFileUsed())
	}

	return v, v.ReadInConfig() == nil
}
