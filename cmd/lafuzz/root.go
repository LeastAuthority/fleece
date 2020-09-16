package main

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"

  "github.com/spf13/viper"
)


var cfgFile string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "lafuzz",
  Short: "A tool to manage fuzzing environment",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $(pwd)/.lafuzz.yaml)")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
  	pwd, err := os.Getwd()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in pwd directory with name ".cmd" (without extension).
    viper.AddConfigPath(pwd)
  	viper.SetConfigType("yaml")
    viper.SetConfigName(".lafuzz")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

