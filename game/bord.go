package game

import (
    "errors"
    "fmt"
)

const (
    BORD_WIDTH     = 10
    BORD_HEIGHT    = 18
    BRICK_NONE     = 0
    BRICK_INACTIVE = 1
    BRICK_ACTIVE   = 2
)

type Pos struct {
    x   int32
    y   int32
}

type Bord [BORD_HEIGHT][BORD_WIDTH]int32

func (b *Bord) get(x int32, y int32) (v int32, err error) {
    if x >= BORD_WIDTH || y >= BORD_HEIGHT || x < 0 || y < 0 {
        //console.log("Bord get, x or y excced")
        err = errors.New("Bord get, x or y excced")
        return
    }
    return b[y][x], nil
}

func (b *Bord) set(x int32, y int32, v int32) (err error) {
    if x >= BORD_WIDTH || y >= BORD_HEIGHT || x < 0 || y < 0 {
        //console.log("Bord get, x or y excced")
        err = errors.New("Bord set, x or y excced")
        return
    }
    b[y][x] = v
    return nil
}

func (b *Bord) clear() {
    for i := 0; i < BORD_HEIGHT; i++ {
        for j := 0; j < BORD_WIDTH; j++ {
            b[i][j] = BRICK_NONE
        }
    }
}

func (b *Bord) diff(other Bord) (r []Pos) {
    for i := 0; i < BORD_HEIGHT; i++ {
        for j := 0; j < BORD_WIDTH; j++ {
            if b[i][j] != other[i][j] {
                r = append(r, Pos{x: int32(j), y: int32(i)})
            }
        }
    }
    return
}

func (b *Bord) diff2(other Bord) (r map[Pos]int32) {
    m := make(map[Pos]int32, 0)
    for i := 0; i < BORD_HEIGHT; i++ {
        for j := 0; j < BORD_WIDTH; j++ {
            if b[i][j] != other[i][j] {
                //r = append(r, Pos{x: int32(j), y: int32(i)})
                pos := Pos{x: int32(j), y: int32(i)}
                m[pos] = b[i][j]
            }
        }
    }
    return m

}

func (b *Bord) bordEnd() bool {
    for i := 0; i < BORD_WIDTH; i++ {
        if b[0][i] == BRICK_INACTIVE {
            return true
        }
    }
    return false
}

func (b *Bord) clearFullLine(line int32) {
    for i := line; i > 0; i-- {
        for j := 0; j < BORD_WIDTH; j++ {
            b[i][j] = b[i-1][j]
        }
    }

    for i := 0; i < BORD_WIDTH; i++ {
        b[0][i] = 0
    }
}

func (b *Bord) getFullLine() (lines []int32) {
    for i := 0; i < BORD_HEIGHT; i++ {
        full := true
        for j := 0; j < BORD_WIDTH; j++ {
            if b[i][j] != BRICK_INACTIVE {
                full = false
                break
            }
        }
        if full {
            lines = append(lines, int32(i))
        }
    }
    fmt.Printf("getFullLine, lines=%v\n", lines)
    return
}

func (b *Bord) clearFullLines() {
    fullLines := b.getFullLine()

    for _, line := range fullLines {
        b.clearFullLine(line)
    }

    return

}

// func main() {
//     var b1 = Bord{}
//     b1.set(3, 3, 5)
//     fmt.Printf("%v", b1)

//     var b2 = Bord{}
//     b2.set(3, 3, 5)
//     r := b2.diff(b1)

//     fmt.Printf("r= %v", r)
// }
