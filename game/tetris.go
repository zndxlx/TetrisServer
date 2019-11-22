package game

import (
    "fmt"
)

type (
    Tetris struct {
        uid  int64
        bord Bord
        // oldbord           Bord //用于判断需要更新的数据
        shape             Shape
        shapeX            int32
        shapeY            int32
        lastActivePosList []Pos
        bEnd              bool
        uFunc             UpdateFunc
    }
)

func NewTetris(uid int64, fn UpdateFunc) *Tetris {
    t := &Tetris{uid: uid, lastActivePosList: []Pos{}, bEnd: false, uFunc: fn}
    return t
}

func (t *Tetris) init() {
    t.generateBrick()
}

func (t *Tetris) reset() {
    t.shape = newRandomShape()
    t.shapeX = 4
    t.shapeY = -4
    t.bEnd = false
    t.lastActivePosList = t.lastActivePosList[0:0]
    t.bord.clear()
    t.update()
}

func (t *Tetris) getActivePosList(px, py int32) (posList []Pos) {
    shape := t.shape.getShape()
    for i := 0; i < len(shape); i++ {
        for j := 0; j < len(shape[0]); j++ {
            //console.log(i, ":", j, this._shape[i][j]);
            if shape[i][j] > 0 {
                posX := px + int32(j)
                posY := py + int32(i)
                posList = append(posList, Pos{posX, posY})
            }
        }
    }
    fmt.Printf("uid=%d, posList=%v\n", t.uid, posList)
    return posList
}

func (t *Tetris) canMove(posX int32, posY int32) bool {
    //1、获取所有需要显示的点
    posList := t.getActivePosList(posX, posY)
    //2、判断每个点的位置是否合法

    for i := 0; i < len(posList); i++ {
        if posList[i].x < 0 || posList[i].x >= BORD_WIDTH {
            return false
        }
        if posList[i].y >= BORD_HEIGHT {
            return false
        }
        v, _ := t.bord.get(posList[i].x, posList[i].y)
        if v == BRICK_INACTIVE {
            //console.log(t.bord.get(posList[i].x, posList[i].y));
            return false
        }
    }

    return true
}

func (t *Tetris) calEx(posX, posY int32) (exLeft, exRight, exBottom int32) {
    //1、获取所有需要显示的点
    posList := t.getActivePosList(posX, posY)
    for i := 0; i < len(posList); i++ {
        if posList[i].x < 0 {
            tmpLeft := 0 - posList[i].x
            if tmpLeft > exLeft {
                exLeft = tmpLeft
            }
        }

        if posList[i].x >= BORD_WIDTH {
            tmpRight := posList[i].x - BORD_WIDTH + 1
            if tmpRight > exRight {
                exRight = tmpRight
            }
        }

        if posList[i].y >= BORD_HEIGHT {
            tmpBottom := posList[i].y - BORD_HEIGHT + 1
            if tmpBottom > exBottom {
                exBottom = tmpBottom
            }
        }
    }

    return
}

func (t *Tetris) update() {
    //将原有active清空
    for _, pos := range t.lastActivePosList {
        t.bord.set(pos.x, pos.y, 0)
    }
    //console.log(t._lastActivePosList);

    posList := t.getActivePosList(t.shapeX, t.shapeY)
    for _, pos := range posList {
        t.bord.set(pos.x, pos.y, BRICK_ACTIVE)
    }

    //fmt.Printf(this._bord);
    t.uFunc(t.uid, t.bord) //更新页面
    t.lastActivePosList = posList
}

func (t *Tetris) generateBrick() {
    t.lastActivePosList = t.lastActivePosList[0:0]
    //t.shape = newRandomShape()
    t.shape, _ = newShape(1, 0)
    t.shapeX = 4
    t.shapeY = -4
    //posList := t.getActivePosList(t.shapeX, t.shapeY)
    // for _, pos := range posList {
    //     t.bord.set(pos.x, pos.y, 0)
    //     // t.render.update(this._bord);
    // }
    //t.lastActivePosList = posList
}

func (t *Tetris) cfBrick() (bEnd bool) {
    //需要改变砖块状态，
    bEnd = false
    posList := t.getActivePosList(t.shapeX, t.shapeY)
    fmt.Printf("1111111 posList=%v\n", posList)
    for _, pos := range posList {
        t.bord.set(pos.x, pos.y, BRICK_INACTIVE)
        //fmt.Printf("2222 set to BRICK_INACTIVE")
    }

    if t.bord.bordEnd() {
        //t.reset();
        //alert("游戏结束");
        fmt.Printf("游戏结束")
        bEnd = true
        t.bEnd = true
        return
    }

    //消除行
    t.bord.clearFullLines()

    //生成新的砖块
    t.generateBrick()
    t.update()
    return
}

func (t *Tetris) doLeft() {
    fmt.Printf("game DoLeft\n")
    can := t.canMove(t.shapeX-1, t.shapeY)
    //console.log(can);

    if !can {
        fmt.Printf("can not left\n")
        return
    }
    t.shapeX--
    t.update()
}

func (t *Tetris) doRight() {
    fmt.Printf("game DoRight\n")
    can := t.canMove(t.shapeX+1, t.shapeY)
    if !can {
        fmt.Printf("can not right\n")
        return
    }
    t.shapeX++
    t.update()
}

func (t *Tetris) doDown() (bEnd bool) {
    fmt.Printf("game down\n")
    bEnd = false
    can := t.canMove(t.shapeX, t.shapeY+1)
    if !can {
        fmt.Printf("can not down\n")
        bEnd = t.cfBrick()
    } else {
        t.shapeY++
        t.update()
    }
    t.bEnd = bEnd
    return
}

func (t *Tetris) doRotate() {
    fmt.Printf("game rotate\n")
    t.shape.rotate()
    exLeft, exRight, exBottom := t.calEx(t.shapeX, t.shapeY)

    fmt.Printf("exLeft=%d exLeft=%d exBottom=%d\n", exLeft, exRight, exBottom)
    t.shapeX = t.shapeX + exLeft - exRight
    t.shapeY = t.shapeY - exBottom

    t.update()
}

//DoOp 游戏命令处理
func (t *Tetris) DoOp(op int) {
    switch op {
    case GAME_OP_LEFT:
        t.doLeft()
    case GAME_OP_Right:
        t.doRight()
    case GAME_OP_DOWN:
        t.doDown()
    case GAME_OP_ROTATE:
        t.doRotate()
    default:
        fmt.Printf("no support op %s\n", op)
    }
}
