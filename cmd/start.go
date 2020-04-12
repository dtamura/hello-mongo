package cmd

import (
	"fmt"
	"os"

	"github.com/dtamura/hello-mongo/server"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 読み込む設定ファイル名
var cfgFile string

// 読み込んだ設定ファイルの構造体
var appConfig server.AppConfig

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Server",
	Long:  "Start Server",
	RunE: func(cmd *cobra.Command, args []string) error {

		return server.Start(appConfig)

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 設定ファイル名をフラグで受け取る
	startCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file name")
}

func initConfig() {
	// Set Default Value
	viper.SetDefault("mode", "debug")
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetConfigType("yaml")
	viper.SetDefault("mongodb.url", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "test")
	viper.SetDefault("ping.url", "http://localhost:8080")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}
		// Search config in home directory with name ".hello-mongo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hello-mongo")
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// 環境変数で上書き
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
