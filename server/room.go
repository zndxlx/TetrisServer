package server

import (
    "TetrisServer/game"
    "TetrisServer/protocol"
    "errors"
    "fmt"
    "github.com/lonng/nano"
    "github.com/lonng/nano/scheduler"
    "github.com/lonng/nano/session"
    "github.com/pborman/uuid"
    log "github.com/sirupsen/logrus"
    "sync/atomic"
)

const (
    ROOM_STATUS_INIT = iota
    ROOM_STATUS_DESTORY
)

//Room 房间结构描述
type Room struct {
    roomNo string // 房间号
    state  int32  // 状态

    creator   int64 // 创建玩家UID
    createdAt int64 // 创建时间
    players   []*Player
    group     *nano.Group // 组播通道
    // die       chan struct{}
    game   *game.Game
    logger *log.Entry
}

func NewRoom(roomNo string) *Room {
    d := &Room{
        state:   ROOM_STATUS_INIT,
        roomNo:  roomNo,
        players: []*Player{},
        group:   nano.NewGroup(uuid.New()),
        //die:     make(chan struct{}),

        logger: log.WithField("room", roomNo),
    }

    return d
}

func (r *Room) members() []int64 {
    uidList := []int64{}
    for _, p := range r.players {
        uidList = append(uidList, p.Uid())
    }
    return uidList
}

// 摧毁桌子
func (r *Room) destroy() {
    if r.status() == ROOM_STATUS_DESTORY {
        r.logger.Info("桌子已经解散")
        return
    }

    // 标记为销毁
    r.setStatus(ROOM_STATUS_DESTORY)

    r.logger.Info("销毁房间")
    for i, p := range r.players {
        r.logger.Debugf("销毁房间，清除玩家%d数据", p.Uid())
        //p.reset()
        p.room = nil
        p.logger = log.WithField("plaer", p.uid)
        r.players[i] = nil
    }

    // 释放room资源
    r.group.Close()

    //删除桌子
    scheduler.PushTask(func() {
        defaultRoomManager.setRoom(r.roomNo, nil)
    })
}

func (r *Room) isDestroy() bool {
    return r.status() == ROOM_STATUS_INIT
}

func (r *Room) setStatus(s int32) {
    atomic.StoreInt32((*int32)(&r.state), int32(s))
}

func (r *Room) status() int32 {
    return atomic.LoadInt32((*int32)(&r.state))
}

func (r *Room) playerJoin(s *session.Session) error {
    uid := s.UID()
    var (
        p *Player
        //err error
    )
    r.group.Add(s)
    exists := false
    for _, p := range r.players {
        if p.Uid() == uid {
            exists = true
            p.logger.Warn("玩家已经在房间中")
            break
        }
    }
    if !exists {
        p = s.Value(CUR_PLAYER).(*Player)
        r.players = append(r.players, p)
        for i, p := range r.players {
            p.setRoom(r, i)
        }
    }
    r.notifyState()
    if r.game == nil && len(r.players) == 2 {
        r.game = game.NewGame(r.players[0].Uid(), r.players[1].Uid(), func(uid int64, b game.Bord) {
            //广播该消息
            res := &protocol.BordStatus{
                Uid:  uid,
                Bord: b,
            }
            r.group.Broadcast("OnGameStateChange", res)
        })
    }
    return nil
}

func (r *Room) onPlayerExit(s *session.Session) {
    uid := s.UID()
    r.group.Leave(s)

    restPlayers := []*Player{}
    for _, p := range r.players {
        if p.Uid() != uid {
            restPlayers = append(restPlayers, p)
        } else { //删除玩家的房间信息
            p.room = nil
            p.turn = 0
        }
    }
    r.players = restPlayers
    r.notifyState()
}

func (r *Room) notifyState() {
    res := &protocol.RoomInfo{
        RoomNo:    r.roomNo,
        CreatedAt: r.createdAt,
        Creator:   r.creator,
        Members:   r.members(),
        Status:    r.status(),
    }

    r.group.Broadcast("OnRoomStateChange", res)
    return
}

func (r *Room) processGameCmd(s *session.Session, req *protocol.GameCmd) error {
    if r.game == nil {
        return errors.New("game not init")
    }
    uid := s.UID()

    switch req.Op {
    case game.GAME_OP_LEFT:
        r.game.DoOp(game.GAME_OP_LEFT, uid)
    case game.GAME_OP_Right:
        r.game.DoOp(game.GAME_OP_Right, uid)
    case game.GAME_OP_DOWN:
        r.game.DoOp(game.GAME_OP_DOWN, uid)
    case game.GAME_OP_ROTATE:
        r.game.DoOp(game.GAME_OP_ROTATE, uid)
    case game.GAME_CMD_PAUSE:
        r.game.DoPause(uid)
    case game.GAME_CMD_START:
        r.game.DoStart(uid)
    default:
        return errors.New("cmd not support")
    }
    return nil
}

var gRoomNo int32 = 0

func newRoomNo() string {
    r := atomic.AddInt32(&gRoomNo, 1)

    return fmt.Sprintf("%d", r)
}
