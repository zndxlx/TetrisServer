class Bord {
    constructor(w, h) {
        this._width = w || 10;
        this._height = h || 18;
        this._bord = new Array(this._height);
        for (let i = 0; i < this._height; i++) {
            this._bord[i] = new Array(this._width);
            for (let j = 0; j < this._width; j++) {
                this._bord[i][j] = 0;
            }
        }
    }
    getW() {
        return this._width;
    }

    getH() {
        return this._height;
    }

    get(x, y) {
        if (x >= this._width || y >= this._height || x < 0 || y < 0) {
            console.log("Bord get, x or y excced");
            return undefined;
        }
        return this._bord[y][x];
    }
    set(x, y, v) {
        if (x >= this._width || y >= this._height || x < 0 || y < 0) {
            console.log("Bord set, x or y excced");
            return undefined;
        }
        this._bord[y][x] = v;
        return this;
    }
    //比较两个border，返回不同的地方坐标列表
    diff(bord) {
        let res = []
        for (let i = 0; i < this._height; i++) {
            for (let j = 0; j < this._width; j++) {
                if (this._bord[i][j] != bord.get(j, i)) {
                    res.push({ x: j, y: i })
                }
            }
        }
        return res
    }

    clear() {
        for (let i = 0; i < this._height; i++) {
            for (let j = 0; j < this._width; j++) {
                this._bord[i][j] = 0;
            }
        }
    }

    //判断是否到顶部
    bordEnd() {
        for (let i = 0; i < this._width; i++) {
            if (this._bord[0][i] == 1) {
                return true;
            }
        }
        return false
    }

    //获取完整行坐标
    getFullLine() {
        let res = []
        for (let i = 0; i < this._height; i++) {
            let full = true;
            for (let j = 0; j < this._width; j++) {
                if (this._bord[i][j] != 1) {
                    full = false;
                    break;
                }
            }
            if (full) {
                res.push(i);
            }
        }
        return res
    }

    clearFullLine(line) {
        console.log(`before line=${line}`);
        console.log(JSON.stringify(this._bord));
        for (let j = line; j > 0; j--) {
            for (let k = 0; k < this._width; k++) {
                this._bord[j][k] = this._bord[j - 1][k];
            }
        }

        for (let i = 0; i < this._width; i++) {
            this._bord[0][i] = 0;
        }
        console.log(`after line=${line}`);
        console.log(JSON.stringify(this._bord));

    }

    clearFullLines() {
        console.log("clearFullLines");

        let fullLines = this.getFullLine();
        console.log(`fullLines=${fullLines}`);

        fullLines.forEach((line) => {
            this.clearFullLine(line);
        })

        return
    }
}


//根据border数据刷新页面显示
class Render {
    constructor() {
        this._oldbord = new (Bord);
        //this._newbord = new (Bord);
        this.BLOCK_SIZE = 20;
        this._mainEl = document.querySelector('.main');
        this.init();
    }

    update(bord) {
        let res = this._oldbord.diff(bord);
        //console.log(res);

        res.forEach((item) => {
            //console.log(item);
            let index = item.y * this._oldbord.getW() + item.x;
            let el = this._mainEl.children[index];
            let v = bord.get(item.x, item.y);
            if (v == 2) {
                el.className = 'brick brickActivity';
            } else if (v == 0) {
                el.className = 'brick brickNone';
            } else if (v == 1) {
                el.className = 'brick brickInactivity';
            } else {
                console.log(`unknow value `);
                //continue;
            }

            this._oldbord.set(item.x, item.y, bord.get(item.x, item.y));
        });

    }

    init() {
        for (let i = 0; i < this._oldbord.getH(); i++) {
            for (let j = 0; j < this._oldbord.getW(); j++) {
                //console.log('hahahah');

                let brick = document.createElement('div');
                brick.className = 'brick brickNone';

                brick.style.top = `${i * this.BLOCK_SIZE}px`;
                brick.style.left = `${j * this.BLOCK_SIZE}px`;
                this._mainEl.appendChild(brick);
            }
        }
    }
}



class ShapeA {
    constructor(t) {   //-1 表示随机生成
        let type = t || 0;
        this._index = 0;
        if (type == -1) {
            type = Math.floor(Math.random() * 7) + 1;
            this.generate(type);
            this._index = Math.floor(Math.random() * this._shapes.length);
        } else {
            this.generate(type);
            this._index = 0;
        }

    }

    generate(type) {
        switch (type) {
            case 0: //空的
                this._shapes = [[], []];
                break;
            case 1: //』
                this._shapes = [[[0, 0, 1], [0, 0, 1], [0, 1, 1]],
                [[0, 0, 0], [1, 0, 0], [1, 1, 1]],
                [[1, 1, 0], [1, 0, 0], [1, 0, 0]],
                [[1, 1, 1], [0, 0, 1], [0, 0, 0]]];
                break;
            case 2: //I
                this._shapes = [[[0, 1, 0, 0], [0, 1, 0, 0], [0, 1, 0, 0], [0, 1, 0, 0]],
                [[0, 0, 0, 0], [1, 1, 1, 1], [0, 0, 0, 0], [0, 0, 0, 0]]];
                break;
            case 3: //田
                this._shapes = [[[1, 1], [1, 1]]];
                break;
            case 4://T
                this._shapes = [[[1, 1, 1], [0, 1, 0], [0, 0, 0]],
                [[0, 0, 1], [0, 1, 1], [0, 0, 1]],
                [[0, 0, 0], [0, 1, 0], [1, 1, 1]],
                [[1, 0, 0], [1, 1, 0], [1, 0, 0]]];
                break;
            case 5://Z
                this._shapes = [[[1, 1, 0], [0, 1, 1], [0, 0, 0]],
                [[0, 1, 0], [1, 1, 0], [1, 0, 0]]];
                break;
            case 6://倒Z
                this._shapes = [[[0, 1, 1], [1, 1, 0]],
                [[1, 0], [1, 1], [0, 1]]];
                break;
            case 7:////L
                this._shapes = [[[1, 0, 0], [1, 0, 0], [1, 1, 0]],
                [[1, 1, 1], [1, 0, 0], [0, 0, 0]],
                [[0, 1, 1], [0, 0, 1], [0, 0, 1]],
                [[0, 0, 0], [0, 0, 1], [1, 1, 1]]];
                break;
            default:
                this._shapes = [[], []];
                break;
        }
    }

    getShape() {
        return this._shapes[this._index];
    }

    rotate() {
        this._index = (this._index + 1) % this._shapes.length;
    }
    // getWidth() {
    //     return this._shapes[this._index][0].length;
    // }
    // getHeight() {
    //     return this._shapes[this._index].length;
    // }
}




class Game {
    constructor() {
        this._bord = new Bord();
        this._render = new Render();
        this._shape = new ShapeA(0);

        this._shape_x = 0;
        this._shape_y = 0;
        this._lastActivePosList = []
        this.init();
    }

    reset() {
        //console.log("hhhh");
        this._shape = new ShapeA(0);
        clearInterval(this.fallDownTimer);
        this.fallDownTimer = null;
        this._shape_x = 0;
        this._shape_y = 0;
        this._lastActivePosList = [];
        this._bord.clear();
        this.update();
    }


    canMove(posX, posY) {
        //1、获取所有需要显示的点
        let posList = this.getActivePosList(posX, posY);
        //2、判断每个点的位置是否合法

        for (let i = 0; i < posList.length; i++) {
            if (posList[i].x < 0 || posList[i].x >= this._bord.getW()) {
                return false;
            }
            if (posList[i].y >= this._bord.getH()) {
                return false;
            }
            if (this._bord.get(posList[i].x, posList[i].y) == 1) {
                console.log(this._bord.get(posList[i].x, posList[i].y));
                return false;
            }
        }

        return true;
    }

    /*计算超出的部分*/
    calEx(posX, posY) {
        let res = {
            exLeft: 0,
            exRight: 0,
            exBottom: 0
        };

        //1、获取所有需要显示的点
        let posList = this.getActivePosList(posX, posY);
        for (let i = 0; i < posList.length; i++) {
            if (posList[i].x < 0) {
                let exLeft = 0 - posList[i].x
                if (exLeft > res.exLeft) {
                    res.exLeft = exLeft;
                }
            }

            if (posList[i].x >= this._bord.getW()) {
                let exRight = posList[i].x - this._bord.getW() + 1;
                if (exRight > res.exRight) {
                    res.exRight = exRight;
                }
            }

            if (posList[i].y >= this._bord.getH()) {
                let exBottom = posList[i].y - this._bord.getH() + 1;
                if (exBottom > res.exBottom) {
                    res.exBottom = exBottom;
                }
            }
        }

        return res;
    }

    left() {
        console.log('game left');
        let can = this.canMove(this._shape_x - 1, this._shape_y);
        console.log(can);

        if (!can) {
            console.log('can not left');
            return
        }
        this._shape_x -= 1;
        this.update()
    }
    right() {
        console.log('game right');
        let can = this.canMove(this._shape_x + 1, this._shape_y);
        if (!can) {
            console.log('can not right');
            return;
        }
        this._shape_x += 1;
        this.update();
    }
    down() {
        console.log('game down');
        let can = this.canMove(this._shape_x, this._shape_y + 1);
        if (!can) {
            console.log('can not down');

            this.cfBrick();

        } else {
            this._shape_y += 1;
            this.update();
        }
    }

    cfBrick() {
        //需要改变砖块状态，
        let posList = this.getActivePosList(this._shape_x, this._shape_y);
        //根据posList更新_bord
        posList.forEach((item) => {
            console.log(item);
            this._bord.set(item.x, item.y, 1);  //更新砖块状态为确认状态
            this._render.update(this._bord);
        });

        if (this._bord.bordEnd()) {
            this.reset();
            alert("游戏结束");
            return;
        }

        //消除行
        this._bord.clearFullLines();
        // this.update();
        //生成新的砖块
        this.generateBrick();
    }

    generateBrick() {
        this._lastActivePosList = [];
        this._shape = new ShapeA(-1);
        this._shape_x = Math.floor(this._bord.getW() / 2) - 1;
        this._shape_y = -2;
        let posList = this.getActivePosList(this._shape_x, this._shape_y);
        //根据posList更新_bord
        posList.forEach((item) => {
            console.log(item);
            this._bord.set(item.x, item.y, 2);
            this._render.update(this._bord);
        });

        this._lastActivePosList = posList;
    }

    rotate() {
        console.log('game rotate');
        this._shape.rotate();
        let exRes = this.calEx(this._shape_x, this._shape_y);
        console.log(exRes);

        this._shape_x = this._shape_x + exRes.exLeft - exRes.exRight;
        this._shape_y = this._shape_y - exRes.exBottom;

        this.update();
    }

    update() {
        //将原有active清空
        this._lastActivePosList.forEach((item) => {
            this._bord.set(item.x, item.y, 0);
        });
        //console.log(this._lastActivePosList);

        let posList = this.getActivePosList(this._shape_x, this._shape_y);
        posList.forEach((item) => {
            // console.log(item);
            this._bord.set(item.x, item.y, 2);   //设置新的acive位置
        });
        //console.log(this._bord);
        this._render.update(this._bord);   //更新页面
        this._lastActivePosList = posList;
    }

    getActivePosList(px, py) {
        let posList = [];
        let shape = this._shape.getShape();
        for (let i = 0; i < shape.length; i++) {
            for (let j = 0; j < shape[0].length; j++) {
                //console.log(i, ":", j, this._shape[i][j]);
                if (shape[i][j] > 0) {
                    let posX = px + j;
                    let posY = py + i;
                    posList.push({ x: posX, y: posY });
                }
            }
        }
        console.log(posList);

        return posList
    }

    init() {
        this.generateBrick();

        this.fallDownTimer = setInterval(() => {
            //console.log("fall Down");
            this.down()
        }, 600);

    }

}


function control(game) {
    document.onkeydown = (e) => {
        e.preventDefault();
        const key = e.keyCode;
        switch (key) {
            case 37:
                console.log('key left');
                game.left();
                break;
            //up
            case 38:
                console.log('key up');
                game.rotate();
                break;
            //right
            case 39:
                console.log('key right');
                game.right();
                break;
            //down
            case 40:
                console.log('key down');
                game.down();
                break;
            default:
                break;
        }
    }
}


// 18 * 10
window.onload = () => {
    // let bord = new Bord();
    // bord.set(5, 4, 10);

    // let rend = new Render();
    // rend.update(bord);
    game = new Game();
    control(game);
}