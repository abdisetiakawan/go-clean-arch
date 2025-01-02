package main

import (
	"fmt"

	"github.com/abdisetiakawan/go-clean-arch/internal/config"
	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
)

func main() {
    viperConfig := config.NewViper()
    log := config.NewLogger(viperConfig)
    db := config.NewDatabase(viperConfig, log)
    validator := config.NewValidator(viperConfig)
    app := config.NewFiber(viperConfig)
    helper := helper.NewHelper(viperConfig)

    config.Bootstrap(&config.BootstrapConfig{
        DB:       db,
        App:      app,
        Log:      log,
        Validate: validator,
        Config:   viperConfig,
        Help:     &helper,
    })

    webPort := viperConfig.GetInt("web.port")
    err := app.Listen(fmt.Sprintf(":%d", webPort))
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}