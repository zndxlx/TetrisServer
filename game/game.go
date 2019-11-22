package game

import (
    "fmt"
    // "github.com/lonng/nano"
    "github.com/lonng/nano/scheduler"
    // log "github.com/sirupsen/logrus"
    "time"
)

const (
    GAME_OP_LEFT   = 1
    GAME_OP_Right  = 2
    GAME_OP_DOWN   = 3
    GAME_OP_ROTATE = 4
    GAME_CMD_PAUSE = 5
    GAME_CMD_START = 6
    SPEED          = 500
)

type (
    UpdateFunc func(uid int64, b Bord)
    InitState  struct{}
    RunState   struct{}
    PauseState struct{}
    OverState  struct{}
    //GameState 游戏状态接口
    GameState interface {
        doOp(game *Game, op int, uid int64)
        doPause(game *Game, uid int64)
        doStart(game *Game, uid int64)
    }
    //Game 游戏
    Game struct {
        tetrises map[int64]*Tetris
        // P1  *Tetris
        // P2  *Tetris
        //Status int32
        timer  *scheduler.Timer
        status GameState
        uFunc  UpdateFunc
    }
)

//NewGame 创建一个游戏
func NewGame(uid1, uid2 int64, uFunc UpdateFunc) *Game {
    // P1 := NewTetris(uid1)
    // P2 := NewTetris(uid2)

    g := &Game{status: InitState{}, tetrises: map[int64]*Tetris{uid1: NewTetris(uid1, uFunc), uid2: NewTetris(uid2, uFunc)}}
    for _, tetris := range g.tetrises {
        tetris.init()
    }
    fmt.Printf("newGame %+v", *g)
    return g
}

//DoOp 游戏控制命令处理
func (g *Game) DoOp(op int, uid int64) {
    g.status.doOp(g, op, uid)
}
func (g *Game) DoPause(uid int64) {
    g.status.doPause(g, uid)
}
func (g *Game) DoStart(uid int64) {
    g.status.doStart(g, uid)
}

//DoOp 游戏控制命令处理
func (s InitState) doOp(g *Game, op int, uid int64) {
    fmt.Printf("state not support this method")
}

//DoPause 游戏暂停请求
func (s InitState) doPause(g *Game, uid int64) {
    fmt.Printf("state not support this method")
}

//DoStart 游戏开始请求
func (s InitState) doStart(g *Game, uid int64) {
    //游戏开始
    if g.timer != nil {
        g.timer.Stop()
    }

    g.timer = scheduler.NewTimer(SPEED*time.Millisecond, func() {
        for _, tetris := range g.tetrises {
            tetris.doDown()
            if tetris.bEnd {
                g.timer.Stop()
                g.timer = nil
                g.status = OverState{}
            }
        }
    })
    g.status = RunState{}
}

//DoOp 游戏控制命令处理
func (s RunState) doOp(g *Game, op int, uid int64) {
    //游戏控制命令处理
    if tetris, ok := g.tetrises[uid]; ok {
        tetris.DoOp(op)
        if tetris.bEnd {
            g.timer.Stop()
            g.timer = nil
            g.status = OverState{}
        }
    } else {
        fmt.Printf("not found uid")
    }

}

//DoPause 游戏暂停请求
func (s RunState) doPause(g *Game, uid int64) {
    //暂停处理
    if g.timer != nil {
        g.timer.Stop()
        g.timer = nil
    }
    g.status = PauseState{}
}

//DoStart 游戏开始请求
func (s RunState) doStart(g *Game, uid int64) {
    fmt.Printf("state not support this method")
}

//DoOp 游戏控制命令处理
func (s PauseState) doOp(g *Game, op int, uid int64) {
    fmt.Printf("state not support this method")
}

//DoPause 游戏暂停请求
func (s PauseState) doPause(g *Game, uid int64) {
    fmt.Printf("state not support this method")
}

//DoStart 游戏开始请求
func (s PauseState) doStart(g *Game, uid int64) {
    //fmt.Printf("state not support this method")
    if g.timer != nil {
        g.timer.Stop()
    }
    g.timer = scheduler.NewTimer(SPEED*time.Millisecond, func() {
        periodDown(g)
    })

    g.status = RunState{}
}

//DoOp 游戏控制命令处理
func (s OverState) doOp(g *Game, op int, uid int64) {
    fmt.Printf("state not support this method")
}

//DoPause 游戏暂停请求
func (s OverState) doPause(g *Game, uid int64) {
    fmt.Printf("state not support this method")
}

//DoStart 游戏开始请求
func (s OverState) doStart(g *Game, uid int64) {
    for _, tetris := range g.tetrises {
        tetris.reset()
    }

    if g.timer != nil {
        g.timer.Stop()
    }

    g.timer = scheduler.NewTimer(SPEED*time.Millisecond, func() {
        periodDown(g)
    })
    g.status = RunState{}
}

func periodDown(g *Game) {
    fmt.Printf("periodDown\n")
    for _, tetris := range g.tetrises {
        tetris.doDown()
        if tetris.bEnd {
            g.timer.Stop()
            g.timer = nil
            g.status = OverState{}
        }
    }
}

// func main() {
//     fmt.Printf("hhhhh\n")
//     g := NewGame(111, 222, func(uid int64, r map[Pos]int32) {
//         fmt.Printf("uid=%d, r=%+v\n", uid, r)
//     })
//     g.Status.DoStart(g, 111)
//     time.Sleep(10 * time.Second)

//     nano.Listen(":3250",
//         nano.WithLogger(log.WithField("component", "nano")),
//     )
// }
