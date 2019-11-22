package server

import (
    "fmt"
    "github.com/lonng/nano"
    "github.com/lonng/nano/component"
    //"github.com/lonng/nano/scheduler"
    "TetrisServer/protocol"
    "github.com/lonng/nano/session"
    log "github.com/sirupsen/logrus"
    //"time"
    // "strconv"
)

type (
    //大厅结构描述
    Hall struct {
        component.Base
        group   *nano.Group       // 广播channel
        players map[int64]*Player // 所有的玩家
    }
)

var defaultHall = NewHall()

//NewHall hall 构造函数
func NewHall() *Hall {
    return &Hall{
        group:   nano.NewGroup("_SYSTEM_MESSAGE_BROADCAST"),
        players: map[int64]*Player{},
    }
}

//AfterInit  hall初始化函数
func (h *Hall) AfterInit() {
    session.Lifetime.OnClosed(func(s *session.Session) {
        log.Infof("111111玩家: %d离开", s.UID())
        h.group.Leave(s)
        delete(h.players, s.UID())
    })

}

//Login 登陆处理
func (h *Hall) Login(s *session.Session, req *protocol.LoginRequest) error {
    uid := req.UID
    s.Bind(uid)

    log.Infof("玩家: %d登录: %+v", uid, req)

    name := fmt.Sprintf("%d_name", uid)
    headURL := fmt.Sprintf("%d_head", uid)

    if p, ok := h.player(uid); !ok {
        log.Infof("玩家: %d不在线，创建新的玩家", uid)

        p = newPlayer(s, uid, name, headURL, "127.0.0.1", 1)
        h.setPlayer(uid, p)
    } else {
        log.Infof("玩家: %d已经在线", uid)
        // 移除广播频道
        h.group.Leave(s)

        // 重置之前的session
        if prevSession := p.session; prevSession != nil && prevSession != s {
            // 如果之前房间存在，则退出来
            if p, err := playerWithSession(prevSession); err == nil && p != nil && p.room != nil && p.room.group != nil {
                p.room.group.Leave(prevSession)
            }

            prevSession.Clear()
            prevSession.Close()
        }

        // 绑定新session
        p.bindSession(s)
    }

    // 添加到广播频道
    h.group.Add(s)

    res := &protocol.LoginResponse{
        Code:    0,
        Name:    name,
        Sex:     1,
        HeadUrl: headURL,
    }

    return s.Response(res)
}

func (h *Hall) P2PMsg(s *session.Session, req *protocol.P2PMsg) error {
    t, err := h.group.Member(req.Tuid)
    if err != nil {
        log.Infof("P2PMsg[%d--->%d] , err=%+v", s.UID(), req.Tuid, err)
        return err
    }

    res := &protocol.OnP2PMsg{
        Fuid:    s.UID(),
        Content: req.Content,
    }

    return t.Push("OnP2PMsg", res)
}

func (h *Hall) player(uid int64) (*Player, bool) {
    p, ok := h.players[uid]

    return p, ok
}

func (h *Hall) setPlayer(uid int64, p *Player) {
    if _, ok := h.players[uid]; ok {
        log.Warnf("玩家已经存在，正在覆盖玩家， UID=%d", uid)
    }
    h.players[uid] = p
}

func (h *Hall) sessionCount() int {
    return len(h.players)
}

// func (h *Hall) offline(uid int64) {
//     delete(h.players, uid)
//     log.Infof("玩家: %d从在线列表中删除, 剩余：%d", uid, len(h.players))
// }
