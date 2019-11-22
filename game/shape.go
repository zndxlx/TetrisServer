package game

import (
    "errors"
    //"fmt"
    "math/rand"
)

type ShapeArry [][]int

type ShapeArryList []ShapeArry

var (
    SHAPEA = ShapeArryList{{{0, 0, 1}, {0, 0, 1}, {0, 1, 1}},
        {{0, 0, 0}, {1, 0, 0}, {1, 1, 1}},
        {{1, 1, 0}, {1, 0, 0}, {1, 0, 0}},
        {{1, 1, 1}, {0, 0, 1}, {0, 0, 0}}}
    //I
    SHAPEB = ShapeArryList{{{0, 1, 0, 0}, {0, 1, 0, 0}, {0, 1, 0, 0}, {0, 1, 0, 0}},
        {{0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}}}
    //田
    SHAPEC = ShapeArryList{{{1, 1}, {1, 1}}}
    //T
    SHAPED = ShapeArryList{{{1, 1, 1}, {0, 1, 0}, {0, 0, 0}},
        {{0, 0, 1}, {0, 1, 1}, {0, 0, 1}},
        {{0, 0, 0}, {0, 1, 0}, {1, 1, 1}},
        {{1, 0, 0}, {1, 1, 0}, {1, 0, 0}}}
    //Z
    SHAPEE = ShapeArryList{{{1, 1, 0}, {0, 1, 1}, {0, 0, 0}},
        {{0, 1, 0}, {1, 1, 0}, {1, 0, 0}}}
    //倒Z
    SHAPEF = ShapeArryList{{{0, 1, 1}, {1, 1, 0}},
        {{1, 0}, {1, 1}, {0, 1}}}
    //L
    SHAPEG = ShapeArryList{{{1, 0, 0}, {1, 0, 0}, {1, 1, 0}},
        {{1, 1, 1}, {1, 0, 0}, {0, 0, 0}},
        {{0, 1, 1}, {0, 0, 1}, {0, 0, 1}},
        {{0, 0, 0}, {0, 0, 1}, {1, 1, 1}}}

    SHAPELIST = [...]ShapeArryList{SHAPEA, SHAPEB, SHAPEC, SHAPED, SHAPEE, SHAPEF, SHAPEG}
)

type Shape struct {
    index  int
    shapes ShapeArryList
}

func newRandomShape() Shape {
    shapeTotal := len(SHAPELIST)
    shape := SHAPELIST[rand.Intn(shapeTotal)]
    shapeArrayTotal := len(shape)
    index := rand.Intn(shapeArrayTotal)
    return Shape{index, shape}
}

func newShape(t int, index int) (Shape, error) {
    switch t {
    case 1:
        return Shape{index, SHAPEA}, nil
    case 2:
        return Shape{index, SHAPEB}, nil
    case 3:
        return Shape{index, SHAPEC}, nil
    case 4:
        return Shape{index, SHAPED}, nil
    case 5:
        return Shape{index, SHAPEE}, nil
    case 6:
        return Shape{index, SHAPEF}, nil
    case 7:
        return Shape{index, SHAPEG}, nil
    default:
        return Shape{}, errors.New("not support shape type")
    }
    return Shape{}, errors.New("not support shape type")
}

func (s *Shape) getShape() [][]int {
    return s.shapes[s.index]
}

func (s *Shape) rotate() {
    s.index = (s.index + 1) % len(s.shapes)
}

// func main() {
//     //s, _ := newShape(1, 0)
//     //fmt.Printf("SHAPEA=%+v", SHAPEA)
//     s := newRandomShape()
//     fmt.Printf("s=%+v\n", s)
//     fmt.Printf("c=%+v\n", s.getShape())

// }
