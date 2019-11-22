package server

import (
    //"fmt"
    "TetrisServer/protocol"
    "errors"
    "github.com/lonng/nano/component"
    "github.com/lonng/nano/scheduler"
    "github.com/lonng/nano/session"
    "time"
)

type (
    //RoomManager 房间管理
    RoomManager struct {
        component.Base
        rooms map[string]*Room
    }
)

var defaultRoomManager = NewRoomManager()

//NewRoomManager 构造器
func NewRoomManager() *RoomManager {
    return &RoomManager{
        rooms: map[string]*Room{},
    }
}

//AfterInit 初始化处理
func (manager *RoomManager) AfterInit() {
    session.Lifetime.OnClosed(func(s *session.Session) {
        if err := manager.onPlayerDisconnect(s); err != nil {
            logger.Errorf("玩家退出: UID=%d, Error=%s", s.UID, err.Error())
        }

    })

    // 每5分钟清空一次已摧毁的房间信息
    scheduler.NewTimer(300*time.Second, func() {
        destroyRooms := map[string]*Room{}
        deadline := time.Now().Add(-24 * time.Hour).Unix()
        for no, d := range manager.rooms {
            // 清除创建超过24小时的房间
            if d.status() == ROOM_STATUS_DESTORY || d.createdAt < deadline {
                destroyRooms[no] = d
            }
        }
        for _, r := range destroyRooms {
            r.destroy()
        }
    })
}

func (manager *RoomManager) onPlayerDisconnect(s *session.Session) error {
    //uid := s.UID()
    p, err := playerWithSession(s)
    if err != nil {
        return err
    }
    p.logger.Debug("RoomManager.onPlayerDisconnect: 玩家网络断开")

    // 移除session
    p.removeSession()

    room := p.room
    room.onPlayerExit(s)
    return nil
}

func (manager *RoomManager) room(number string) (*Room, bool) {
    r, ok := manager.rooms[number]
    return r, ok
}

func (manager *RoomManager) setRoom(number string, room *Room) {
    if room == nil {
        delete(manager.rooms, number)
        logger.WithField("room", number).Debugf("清除房间(%s): 剩余: %d", number, len(manager.rooms))
    } else {
        manager.rooms[number] = room
    }
}

func (manager *RoomManager) CreateRoom(s *session.Session, req *protocol.CreateRoomRequest) error {
    p, err := playerWithSession(s)
    if err != nil {
        return err
    }

    if p.room != nil {
        //重复进入房间
        res := &protocol.CreateRoomResponse{
            Code: 1,
        }

        return s.Response(res)
    }

    no := newRoomNo()
    r := NewRoom(no)
    r.createdAt = time.Now().Unix()
    r.creator = s.UID()
    //房间创建者自动join
    if err := r.playerJoin(s); err != nil {
        return nil
    }

    manager.rooms[no] = r

    res := &protocol.CreateRoomResponse{
        Code: 0,
        Info: protocol.RoomInfo{
            RoomNo:    no,
            CreatedAt: r.createdAt,
            Creator:   r.creator,
            Status:    r.state,
            Members:   []int64{r.creator},
        },
    }

    return s.Response(res)
}

func (manager *RoomManager) JoinRoom(s *session.Session, req *protocol.JoinRoomRequest) error {
    r, ok := manager.room(req.RoomNo)
    if !ok {
        res := &protocol.JoinRoomResponse{
            Code:  1,
            Error: "房间不存在",
        }
        return s.Response(res)
    }

    if err := r.playerJoin(s); err != nil {
        r.logger.Errorf("玩家加入房间失败，UID=%d, Error=%s", s.UID(), err.Error())
        res := &protocol.JoinRoomResponse{
            Code:  2,
            Error: err.Error(),
        }
        return s.Response(res)
    }

    res := &protocol.JoinRoomResponse{
        Code: 0,
        Info: protocol.RoomInfo{
            RoomNo:    r.roomNo,
            CreatedAt: r.createdAt,
            Creator:   r.creator,
            Members:   r.members(),
            Status:    r.status(),
        },
    }

    return s.Response(res)
}

func (manager *RoomManager) RoomMsg(s *session.Session, req *protocol.RoomMsg) error {
    p, err := playerWithSession(s)
    if err != nil {
        return err
    }

    if p.room == nil {
        //重复进入房间
        res := &protocol.OnP2PMsg{
            Fuid:    1000,
            Content: "没有进入房间",
        }

        return s.Push("OnP2PMsg", res)
    }

    res := &protocol.OnRoomMsg{
        Fuid:    s.UID(),
        RoomID:  p.room.roomNo,
        Content: req.Content,
    }

    return p.room.group.Broadcast("OnRoomMsg", res)
}

func (manager *RoomManager) GameCmd(s *session.Session, req *protocol.GameCmd) error {
    p, err := playerWithSession(s)
    if err != nil {
        return err
    }

    if p.room == nil {
        //重复进入房间
        return errors.New("need enter room")
    }
    logger.Debugf("GameCmd req=%v", *req)
    return p.room.processGameCmd(s, req)
}
