package server

import (
    "github.com/lonng/nano/session"
    log "github.com/sirupsen/logrus"
)

const (
    //CUR_PLAYER 绑定player到session中的key
    CUR_PLAYER = "kCurPlayer"
)

//Player 玩家结构描述
type Player struct {
    uid  int64  // 用户ID
    head string // 头像地址
    name string // 玩家名字
    ip   string // ip地址
    sex  int    // 性别

    // 玩家数据
    session *session.Session

    // 游戏相关字段

    room   *Room      //当前房间
    turn   int        //在房间的位置
    logger *log.Entry // 日志
}

func newPlayer(s *session.Session, uid int64, name, head, ip string, sex int) *Player {
    p := &Player{
        uid:  uid,
        name: name,
        head: head,
        ip:   ip,
        sex:  sex,

        logger: log.WithField("player", uid),
    }
    p.bindSession(s)
    return p
}

func (p *Player) bindSession(s *session.Session) {
    p.session = s
    p.session.Set(CUR_PLAYER, p)
}

func (p *Player) removeSession() {
    p.session.Remove(CUR_PLAYER)
    p.session = nil
}

//Uid 返回uid
func (p *Player) Uid() int64 {
    return p.uid
}

func (p *Player) setRoom(r *Room, turn int) {
    if r == nil {
        p.logger.Error("桌号为空")
        return
    }

    p.room = r
    p.turn = turn

    p.logger = log.WithFields(log.Fields{"room": p.room.roomNo, "player": p.uid})

}

func playerWithSession(s *session.Session) (*Player, error) {
    p, ok := s.Value(CUR_PLAYER).(*Player)
    if !ok {
        return nil, ErrPlayerNotFound
    }
    return p, nil
}
