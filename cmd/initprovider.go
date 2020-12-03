package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/drakkan/sftpgo/config"
	"github.com/drakkan/sftpgo/dataprovider"
	"github.com/drakkan/sftpgo/logger"
	"github.com/drakkan/sftpgo/utils"
)

var (
	initProviderCmd = &cobra.Command{
		Use:   "initprovider",
		Short: "Initializes and/or updates the configured data provider",
		Long: `This command reads the data provider connection details from the specified
configuration file and creates the initial structure or update the existing one,
as needed.

Some data providers such as bolt and memory does not require an initialization
but they could require an update to the existing data after upgrading SFTPGo.

For SQLite/bolt providers the database file will be auto-created if missing.

For PostgreSQL and MySQL providers you need to create the configured database,
this command will create/update the required tables as needed.

To initialize/update the data provider from the configuration directory simply use:

$ sftpgo initprovider

Please take a look at the usage below to customize the options.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.DisableLogger()
			logger.EnableConsoleLogger(zerolog.DebugLevel)
			configDir = utils.CleanDirInput(configDir)
			err := config.LoadConfig(configDir, configFile, viperInstance)
			if err != nil {
				logger.WarnToConsole("Unable to initialize data provider, config load error: %v", err)
				return
			}
			providerConf := config.GetProviderConf()
			logger.InfoToConsole("Initializing provider: %#v config file: %#v", providerConf.Driver, viperInstance.ConfigFileUsed())
			err = dataprovider.InitializeDatabase(providerConf, configDir)
			if err == nil {
				logger.InfoToConsole("Data provider successfully initialized/updated")
			} else if err == dataprovider.ErrNoInitRequired {
				logger.InfoToConsole("%v", err.Error())
			} else {
				logger.WarnToConsole("Unable to initialize/update the data provider: %v", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initProviderCmd)
	addConfigFlags(initProviderCmd)
}
