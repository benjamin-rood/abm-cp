function DrawList(obj) {
  this.cpp = obj.data.cpp
  this.vp = obj.data.vp
  this.bg = obj.data.bg
}

var initDrawObj = {
  type: 'drawlist',
  data: {
    cpp: [],
    vp: [],
    bg: {
      red: 0,
      green: 0,
      blue: 0
    }
  }
}

var drawlist = new DrawList(initDrawObj)

var sketch = function(p) {
  var modelDw = 1.0
  var modelDh = 1.0
  var vpSize = 30
  p.setup = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.75
    p.createCanvas(w, h)
    p.noLoop()
    p.background(0, 0, 255)
  }

  p.draw = function() {
    p.background(drawlist.bg.red, drawlist.bg.green, drawlist.bg.blue)

    if (drawlist.cpp) {
      for (var i = 0; i < drawlist.cpp.length; i++) {
        var x = absToView(drawlist.cpp[i].position.x, modelDw, p.width)
        var y = absToView(drawlist.cpp[i].position.y, modelDh, p.height)
        var col = p.color(drawlist.cpp[i].colour.red, drawlist.cpp[i].colour.green, drawlist.cpp[i].colour.blue)
        p.strokeWeight(8)
        p.stroke(col)
        p.point(x, y)
      }
    }

    if (drawlist.vp) {
      for (var i = 0; i < drawlist.vp.length; i++) {
        var x = drawlist.vp[i].position.x
        var y = drawlist.vp[i].position.y
        var angle = 0 //  because positive rotations here happen clockwise, rather than the convention of the unit circle
        var col = p.color(drawlist.vp[i].colour.red, drawlist.vp[i].colour.green, drawlist.vp[i].colour.blue)
        p.push()
        p.translate(x, y)
        p.rotate(angle)
        p.fill(col)
        p.strokeWeight(1)
        p.stroke(255)
        p.ellipse(0, 0, vpSize, vpSize)
          // p.triangle(0-tSize, 0+tSize, 0, 0-(2*tSize), 0+tSize, 0+tSize)
        p.pop()
      }
    }
  }

  p.windowResized = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.75
    p.resizeCanvas(w, h)
  }
}

function absToView(p, d, n) {
  view = (((p + d) / (2 * d)) * n)
  return view
}



var wsUrl = "ws://" + window.location.hostname + ":" + window.location.port + "/ws"
console.log(wsUrl)
var vizSocket = new WebSocket(wsUrl)
var viz = new p5(sketch, 'abm-viewport')

vizSocket.onopen = function(e) {
  console.log('viz Stream (WebSocket) is opened.')
}

vizSocket.onmessage = function(e) {
  var rawmsg = JSON.parse(e.data)
    // console.log(rawmsg)
  switch (rawmsg.type) {
    case 'render':
      drawlist.cpp = rawmsg.data.cpp
      drawlist.vp = rawmsg.data.vp
      drawlist.bg = rawmsg.data.bg
      viz.redraw()
      break
    case 'statistics':
      // do something
      console.log("recived statistics")
      break
    default:
      console.log("Error: don't recognise the received JSON message type!")
  }
}

vizSocket.onerror = function(e) {
  console.log("WebSocket Error:" +  e)
  if (e.readyState === vizSocket.CLOSED) {
    alert("WebSocket connection expired.")
    vizSocket.close()
  }
}

$(function () {
  $('#contextParamsSend').on('click', function() {
    var context = {
      ['abm-cpp-pop-start']: parseInt($('#abm-cpp-pop-start').val()),
      ['abm-cpp-pop-cap']: parseInt($('#abm-cpp-pop-cap').val()),
      ['abm-cpp-ageing']: parseBool($('#abm-cpp-ageing').val()),
      ['abm-cpp-lifespan']: parseInt($('#abm-cpp-lifespan').val()),
      ['abm-cpp-speed']: parseInt($('#abm-cpp-speed').val()),
      ['abm-cpp-turn']: parseInt($('#abm-cpp-turn').val()),
      ['abm-cpp-sexual-cost']: parseInt($('#abm-cpp-sexual-cost').val()),
      ['abm-cpp-reproduction-chance']: parseFloat($('#abm-cpp-reproduction-chance').val()),
      ['abm-cpp-gestation']: parseInt($('#abm-cpp-sexual-cost').val()),
      ['abm-cpp-spawn-size']: parseInt($('#abm-cpp-spawn-size').val()),
      ['abm-cpp-mf']: parseFloat($('#abm-cpp-mf').val()),
      ['abm-vp-pop-start']: parseInt($('#abm-vp-pop-start').val()),
      ['abm-vp-pop-cap']: parseInt($('#abm-vp-pop-cap').val()),
      ['abm-vp-ageing']: parseBool($('#abm-vp-ageing').val()),
      ['abm-vp-lifespan']: parseInt($('#abm-vp-lifespan').val()),
      ['abm-vp-speed']: parseFloat($('#abm-vp-speed').val()),
      ['abm-vp-turn']: parseFloat($('#abm-vp-turn').val()),
      ['abm-vp-vsr']: parseFloat($('#abm-vp-vsr').val()),
      ['abm-vp-vsr-chance']: parseFloat($('#abm-vp-vsr-chance').val()),
      ['abm-vp-attack-chance']: parseFloat($('#abm-vp-attack-chance').val()),
      ['abm-vp-col-imprinting']: parseFloat($('#abm-vp-col-imprinting').val()),
      ['abm-random-ages']: parseBool($('#abm-random-ages').val()),
      ['abm-rng-random-seed']: parseBool($('#abm-rng-random-seed').val()),
      ['abm-rng-seedval']: parseInt($('#abm-rng-seedval').val())
    }

    var OutMsg = {
      type: "context",
      data: context
    }

    var json = JSON.stringify(OutMsg, null, 2)
    console.log(json)
    vizSocket.send(json)
  })
})

function parseBool(value){
    if (typeof(value) == 'string'){
        value = value.toLowerCase().trim()
    }
    switch(value){
        case true:
        case "true":
        case 1:
        case "1":
        case "on":
        case "yes":
            return true
        default:
            return false
    }
}
