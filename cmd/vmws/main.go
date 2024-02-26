package main

import (
	"github.com/catundercar/vmws-go/cmd/vmws/internal/vm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vmws",
	Short: "A command line tool to manage vms by wrapped the commands: vmrest, vmrun.",
}

func init() {
	// 加载运行初始化配置
	cobra.OnInitialize(initConfig)
	// rootCmd，命令行下读取配置文件，持久化的 flag，全局的配置文件
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vmws.yaml)")
	// local flag，本地化的配置
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(vm.VMCmd)
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

// 初始化配置的一些设置
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile) // viper 设置配置文件
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".vmws")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil { // 读取配置文件
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		log.Fatal(err)
	}
}
