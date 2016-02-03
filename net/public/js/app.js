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
  var vpSize = 10
  var cpSize = 4
  p.setup = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.65
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
        p.strokeWeight(cpSize)
        p.stroke(col)
        p.point(x, y)
      }
    }

    if (drawlist.vp) {
      for (var i = 0; i < drawlist.vp.length; i++) {
        var x = absToView(drawlist.vp[i].position.x, modelDw, p.width)
        var y = absToView(drawlist.vp[i].position.y, modelDh, p.height)
        var angle = drawlist.vp[i].heading
        var col = p.color(drawlist.vp[i].colour.red, drawlist.vp[i].colour.green, drawlist.vp[i].colour.blue)
        p.fill(col)
        p.noStroke()
        p.push()
          p.translate(x, y)
          p.rotate(p.atan2(1, 0))
          p.rotate(angle)
          p.triangle(-vpSize, vpSize, 0, -vpSize, vpSize, vpSize)
          p.fill(255)
          p.triangle(-vpSize/2, 0, 0, -vpSize, vpSize/2, 0)
        p.pop()
      }
    }
  }

  p.windowResized = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.65
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
  console.dir(rawmsg)
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
      ['abm-cpp-ageing']: parseBool($('#abm-cpp-ageing').is(':checked')),
      ['abm-cpp-lifespan']: parseInt($('#abm-cpp-lifespan').val()),
      ['abm-cpp-speed']: parseFloat($('#abm-cpp-speed').val()),
      ['abm-cpp-turn']: parseFloat($('#abm-cpp-turn').val()),
      ['abm-cpp-sexual-cost']: parseInt($('#abm-cpp-sexual-cost').val()),
      ['abm-cpp-reproduction-chance']: parseFloat($('#abm-cpp-reproduction-chance').val()),
      ['abm-cpp-gestation']: parseInt($('#abm-cpp-sexual-cost').val()),
      ['abm-cpp-spawn-size']: parseInt($('#abm-cpp-spawn-size').val()),
      ['abm-cpp-mf']: parseFloat($('#abm-cpp-mf').val()),
      ['abm-vp-pop-start']: parseInt($('#abm-vp-pop-start').val()),
      ['abm-vp-pop-cap']: parseInt($('#abm-vp-pop-cap').val()),
      ['abm-vp-ageing']: parseBool($('#abm-vp-ageing').is(':checked')),
      ['abm-vp-lifespan']: parseInt($('#abm-vp-lifespan').val()),
      ['abm-starvation']: parseBool($('#abm-starvation').is(':checked')),
      ['abm-vp-starvation-point']: parseInt($('#abm-vp-starvation-point').val()),
      ['abm-vp-speed']: parseFloat($('#abm-vp-speed').val()),
      ['abm-vp-turn']: parseFloat($('#abm-vp-turn').val()),
      ['abm-vp-vsr']: parseFloat($('#abm-vp-vsr').val()),
      ['abm-vp-visual-acuity']: parseFloat($('#abm-vp-visual-acuity').val()),
      ['abm-vp-visual-acuity-bump']: parseFloat($('#abm-vp-visual-acuity-bump').val()),
      ['abm-vp-vsr-chance']: parseFloat($('#abm-vp-vsr-chance').val()),
      ['abm-vp-attack-chance']: parseFloat($('#abm-vp-attack-chance').val()),
      ['abm-vp-col-imprinting']: parseFloat($('#abm-vp-col-imprinting').val()),
      ['abm-vp-reproduction-chance']: parseFloat($('#abm-vp-reproduction-chance').val()),
      ['abm-vp-spawn-size']: parseInt($('#abm-vp-spawn-size').val()),
      ['abm-random-ages']: parseBool($('#abm-random-ages').is(':checked')),
      ['abm-rng-random-seed']: parseBool($('#abm-rng-random-seed').is(':checked')),
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
