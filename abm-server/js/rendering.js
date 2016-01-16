function DrawList (obj) {
  this.cpp = obj.data.cpp
  this.vp = obj.data.vp
  this.bg = obj.data.bg
}

var initDrawObj = {type: 'drawlist', data: {cpp: [], vp: [], bg: {red: 0, green: 0, blue: 0}}}
var drawlist = new DrawList(initDrawObj)
var modelDw = 1.0
var modelDh = 1.0

var sketch = function (p) {
  var vSize = 40

  p.setup = function () {
    p.createCanvas(p.windowWidth, p.windowWidth)
    p.noLoop()
    p.background(0, 0, 255)
  }

  p.draw = function () {
    p.background(drawlist.bg.red, drawlist.bg.green, drawlist.bg.blue)
    for (var i = 0; i < drawlist.cpp.length; i++) {
      var x = absToView(drawlist.cpp[i].position.x, modelDw, p.width)
      var y = absToView(drawlist.cpp[i].position.y, modelDh, p.height)
      var col = p.color(drawlist.cpp[i].colour.red, drawlist.cpp[i].colour.green, drawlist.cpp[i].colour.blue)
      p.strokeWeight(8)
      p.stroke(col)
      p.point(x,y)
    }
    // for (var i = 0; i < drawlist.vp.length; i++) {
    //   var x = drawlist.vp[i].position.x
    //   var y = drawlist.vp[i].position.y
    //   var angle = 0 //  because positive rotations here happen clockwise, rather than the convention of the unit circle
    //   var col = p.color(drawlist.vp[i].colour.red, drawlist.vp[i].colour.green, drawlist.vp[i].colour.blue)
    //   p.push()
    //   p.translate(x, y)
    //   p.rotate(angle)
    //   p.fill(col)
    //   p.strokeWeight(1)
    //   p.stroke(255)
    //   p.ellipse(0, 0, vSize, vSize)
    //   // p.triangle(0-tSize, 0+tSize, 0, 0-(2*tSize), 0+tSize, 0+tSize)
    //   p.pop()
    // }
  }

  p.windowResized = function () {
    p.resizeCanvas(p.windowWidth, p.windowHeight)
  }
}

function absToView(p, d, n) {
  view = (((p + d) / (2 * d)) * n)
  return view
}
