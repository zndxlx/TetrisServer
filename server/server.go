package server

import (
    // "fmt"
    // "math/rand"
    "net/http"
    //"strconv"
    //"strings"
    "time"

    "github.com/lonng/nano"
    "github.com/lonng/nano/component"
    "github.com/lonng/nano/serialize/json"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

var (
    logger = log.WithField("component", "game")
)

//启动入口
func Startup() {
    heartbeat := viper.GetInt("core.heartbeat")
    if heartbeat < 5 {
        heartbeat = 5
    }

    comps := &component.Components{}
    comps.Register(defaultHall)
    comps.Register(defaultRoomManager)

    //log.SetFlags(log.LstdFlags | log.Llongfile)
    http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

    nano.Listen(":3250",
        nano.WithIsWebsocket(true),
        nano.WithTimerPrecision(time.Millisecond*50),
        nano.WithHeartbeatInterval(time.Duration(heartbeat)*time.Second),
        nano.WithCheckOriginFunc(func(_ *http.Request) bool { return true }),
        nano.WithWSPath("/nano"),
        nano.WithDebugMode(),
        nano.WithSerializer(json.NewSerializer()), // override default serializer
        nano.WithComponents(comps),
        nano.WithLogger(log.WithField("component", "nano")),
    )
}
