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
    jwt := helper.NewJWTHelper(viperConfig)
    redisClient := config.NewRedisClient(viperConfig, log)
    cache := helper.NewCacheHelper(redisClient)

    config.Bootstrap(&config.BootstrapConfig{
        DB:       db,
        App:      app,
        Log:      log,
        Validate: validator,
        Config:   viperConfig,
        Jwt:      jwt,
        Cache:    cache,
    })

    webPort := viperConfig.GetInt("web.port")
    err := app.Listen(fmt.Sprintf(":%d", webPort))
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}