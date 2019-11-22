package main

import (
    // "fmt"
    //"net/http"
    // "strings"
    // "time"

    "TetrisServer/server"
    // "github.com/lonng/nano"
    // "github.com/lonng/nano/component"
    // "github.com/lonng/nano/pipeline"
    // "github.com/lonng/nano/scheduler"
    // "github.com/lonng/nano/serialize/json"
    // "github.com/lonng/nano/session"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

func main() {
    //初始化配置
    viper.SetConfigType("toml")
    viper.SetConfigFile("./config/config.toml")
    viper.ReadInConfig()

    //初始化日志
    log.SetFormatter(&log.TextFormatter{DisableColors: true})
    if viper.GetBool("core.debug") {
        log.SetLevel(log.DebugLevel)
    }

    // log.WithFields(log.Fields{
    //     "animal": "walrus",
    // }).Info("%+v", viper.GetBool("core.debug"))

    server.Startup()

}
