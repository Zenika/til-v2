package configuration

import (
	"github.com/spf13/viper"
	"github.com/zenika/tilv2back/internal/structures"
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger
var Configuration = structures.Configuration{
	Debug:            false,
	DatabaseFileName: "til.sqlite3",
	ServerPort:       8000,
	Google: structures.ConfigurationGoogle{
		TokenEndpoint: "https://oauth2.googleapis.com/token",
	},
	UseEmbeddedFrontend: false,
}

func init() {
	// Initialize logger
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("TIL")
	viper.AutomaticEnv()

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		Logger.Error("Error unmarshalling configuration", err)
	}

	if Configuration.Debug {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		Logger.Debug("Debug mode activated! Please be careful, as this mode MAY leak information in logs, and flood your drive.")
	}
}
