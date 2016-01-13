function DrawList (obj) {
  this.cpp = obj.data.cpp
  this.vp = obj.data.vp
  this.bg = obj.data.bg
}

var initDrawObj = {type: 'drawlist', data: {cpp: [], vp: [], bg: {red: 0, green: 0, blue: 0}}}
var drawlist = new DrawList(initDrawObj)
var newWidth = 1440
var newHeight = 900

var sketch = function (p) {
  p.setup = function () {
    p.createCanvas(newWidth, newHeight)
    p.strokeWeight(3)
    // p.stroke(255, 255, 255)
    p.noLoop()
    p.background(0, 0, 255)
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
  }
}

var wsocket = new WebSocket('ws://localhost:8080/ws')

var viz = new p5(sketch, 'abm-viewport')

wsocket.onopen = function () {
  console.log('websocket connection init')
}

wsocket.onmessage = function (e) {
  var rawmsg = JSON.parse(e.data)
  console.log(rawmsg)
  if (rawmsg.type === 'render') {
    drawlist.cpp = rawmsg.data.cpp
    drawlist.vp = rawmsg.data.vp
    drawlist.bg = rawmsg.data.bg
  }
  viz.redraw()
}

wsocket.onerror = function (e) {
  console.log('ERROR: ' + e)
}
