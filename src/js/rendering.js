function DrawList(obj) {
  this.cpp = obj.data.cpp
  this.vp = obj.data.vp
  this.bg = obj.data.bg
}

var jsonValue = JSON.parse(jsonString)
console.log(jsonValue)
var drawlist = new DrawList(jsonValue)


var sketch = function (p) {
  p.setup = function () {
    p.createCanvas(300, 200)
    p.background(0.0, 0.0, 0.0)
    p.noCursor()
    p.frameRate(60)
    // p.noLoop()
  }

  p.draw = function () {
    for (i = 0; i < drawlist.cpp.length; i++) {
      x = drawlist.cpp[i].position.x
      y = drawlist.cpp[i].position.y
      col = drawlist.cpp[i].colour
      p.stroke(col.red, col.green, col.blue, 1.0)
      p.strokeWeight(5.0)
      p.point(x,y)
    }
  }
}

var viz = new p5(sketch, 'abm-viewport')
