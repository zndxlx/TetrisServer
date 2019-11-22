package protocol

import (
    "TetrisServer/game"
)

type LoginRequest struct {
    UID int64 `json:"uid"`
}

type P2PMsg struct {
    Tuid    int64  `json:"tuid"`
    Content string `json:"content"`
}

type OnP2PMsg struct {
    Fuid    int64  `json:"fuid"`
    Content string `json:"content"`
}

type RoomMsg struct {
    Content string `json:"content"`
}

type OnRoomMsg struct {
    Fuid    int64  `json:"fuid"`
    RoomID  string `json:"roomId"`
    Content string `json:"content"`
}

type LoginResponse struct {
    Code    int    `json:"code"`
    Name    string `json:"name"`
    HeadUrl string `json:"headUrl"`
    Sex     int    `json:"sex"`
}

type GameCmd struct {
    Op int `json:"op"`
}

// type ClientGameCmdResponse struct {
//     Code  int    `json:"code"`
//     Error string `json:"error"`
// }

type CreateRoomRequest struct {
    Version string `json:"version"`
}

type RoomInfo struct {
    RoomNo    string  `json:"roomNo"`
    CreatedAt int64   `json:"createdAt"`
    Creator   int64   `json:"creator"`
    Status    int32   `json:"status"`
    Members   []int64 `json:"members"`
}

type CreateRoomResponse struct {
    Code  int      `json:"code"`
    Error string   `json:"error"`
    Info  RoomInfo `json:"info"`
}

type JoinRoomRequest struct {
    Version string `json:"version"`
    RoomNo  string `json:"roomNo"`
}

type JoinRoomResponse struct {
    Code  int      `json:"code"`
    Error string   `json:"error"`
    Info  RoomInfo `json:"roomInfo"`
}

type InviteUserC struct {
    Tuid   string `json:"tuid"`
    RoomID string `json:"roomId"`
}

type InviteUserS struct {
    FUid   int      `json:"fuid"`
    Fname  string   `json:"fname"`
    RoomID RoomInfo `json:"roomId"`
}

type BordStatus struct {
    Uid  int64     `json:"uid"`
    Bord game.Bord `json:"bord"`
}
