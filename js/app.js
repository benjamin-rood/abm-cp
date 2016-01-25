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
    var context = new Object()
    var abmCppPopStart = document.getElementById("abm-cpp-pop-start").value
    context['abm-cpp-pop-start'] = parseInt($('#abm-cpp-pop-start').val())
    context['abm-cpp-pop-cap'] = parseInt($('#abm-cpp-pop-cap').val())
    context['abm-cpp-ageing'] = parseBool($('#abm-cpp-ageing').val())
    context['abm-cpp-lifespan'] = parseInt($('#abm-cpp-lifespan').val())
    context['abm-cpp-speed'] = parseInt($('#abm-cpp-speed').val())
    context['abm-cpp-turn'] = parseInt($('#abm-cpp-turn').val())
    context['abm-cpp-sexual-cost'] = parseInt($('#abm-cpp-sexual-cost').val())
    context['abm-cpp-reproduction-chance'] = parseFloat($('#abm-cpp-reproduction-chance').val())
    context['abm-cpp-gestation'] = parseInt($('#abm-cpp-sexual-cost').val())
    context['abm-cpp-spawn-size'] = parseInt($('#abm-cpp-spawn-size').val())
    context['abm-cpp-mf'] = parseFloat($('#abm-cpp-mf').val())
    context['abm-vp-pop-start'] = parseInt($('#abm-vp-pop-start').val())
    context['abm-vp-pop-cap'] = parseInt($('#abm-vp-pop-cap').val())
    context['abm-vp-ageing'] = parseBool($('#abm-vp-ageing').val())
    context['abm-vp-lifespan'] = parseInt($('#abm-vp-lifespan').val())
    context['abm-vp-speed'] = parseFloat($('#abm-vp-speed').val())
    context['abm-vp-turn'] = parseFloat($('#abm-vp-turn').val())
    context['abm-vp-vsr'] = parseFloat($('#abm-vp-vsr').val())
    context['abm-vp-vsr-chance'] = parseFloat($('#abm-vp-vsr-chance').val())
    context['abm-vp-attack-chance'] = parseFloat($('#abm-vp-attack-chance').val())
    context['abm-vp-col-imprinting'] = parseFloat($('#abm-vp-col-imprinting').val())
    context['abm-random-ages'] = parseBool($('#abm-random-ages').val())
    context['abm-rng-random-seed'] = parseBool($('#abm-rng-random-seed').val())
    context['abm-rng-seedval'] = parseInt($('#abm-rng-seedval').val())
    var OutMsg = Object
    OutMsg.type = "start"
    OutMsg.data = context
    console.dir(OutMsg)
    console.dir(JSON.stringify(OutMsg))
    vizSocket.send(JSON.stringify(OutMsg))
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
