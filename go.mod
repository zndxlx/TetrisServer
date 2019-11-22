module TetrisServer

go 1.13

require (
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/lonng/nano v0.0.0-00010101000000-000000000000
	github.com/pborman/uuid v1.2.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.5.0
)

replace github.com/lonng/nano => ./nano

replace TetrisServer/game => ./game
