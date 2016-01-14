function DrawList (obj) {
  this.cpp = obj.data.cpp
  this.vp = obj.data.vp
  this.bg = obj.data.bg
}

var initDrawObj = {type: 'drawlist', data: {cpp: [], vp: [], bg: {red: 0, green: 0, blue: 0}}}
var drawlist = new DrawList(initDrawObj)
var newWidth = 600
var newHeight = 400

var sketch = function (p) {
  p.setup = function () {
    p.createCanvas(newWidth, newHeight)
    p.noLoop()
    p.strokeWeight(5)
    // p.stroke(255, 255, 255)
    p.background(0, 0, 255)
    var triangleSize = 20
  }

  p.draw = function () {
    p.background(drawlist.bg.red, drawlist.bg.green, drawlist.bg.blue)
    for (var i = 0; i < drawlist.cpp.length; i++) {
      var x = drawlist.cpp[i].position.x
      var y = drawlist.cpp[i].position.y
      var col = p.color(drawlist.cpp[i].colour.red, drawlist.cpp[i].colour.green, drawlist.cpp[i].colour.blue)
      p.stroke(col)
      p.point(x,y)
      // p.fill(col)
      // p.ellipse(x, y, 15, 15)
    }
    for (var i = 0; i < drawlist.vp.length; i++) {
      var x = drawlist.vp[i].position.x
      var y = drawlist.vp[i].position.y
      var angle = -drawlist.vp[i].heading //  because positive rotations here happen clockwise, rather than the convention of the unit circle
      var col = p.color(drawlist.vp[i].colour.red, drawlist.vp[i].colour.green, drawlist.vp[i].colour.blue)
      translate(x, y)
      rotate(angle)
      fill(col)
      triangle(0-tSize, 0+tSize, 0, 0-(2*tSize), 0+tSize, 0+tSize)
      rotate(-angle)
      translate(-x, -y)
    }
  }
}
