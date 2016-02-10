function DrawList(obj) {
  this.cpp = obj.data.cpp
  this.vp = obj.data.vp
  this.bg = obj.data.bg
  this['cpp-pop-string'] = obj.data['cpp-pop-string']
  this['vp-pop-string'] = obj.data['vp-pop-string']
  this['turncount-string'] = obj.data['turncount-string']
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
    },
    'cpp-pop-string': "",
    'vp-pop-string': "",
    'turncount-string': "",
  }
}

var drawlist = new DrawList(initDrawObj)

console.log(drawlist)

var sketch = function(p) {
  var modelDw = 1.0
  var modelDh = 1.0
  var vpSize = 10
  var cpSize = 4
  p.setup = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.625
    p.createCanvas(w, h)
    p.noLoop()
    p.background(0, 0, 255)
  }

  p.draw = function() {
    p.background(25,18,18)
    //  draw colour polymorphic prey agents
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
    //  draw visual predator agents
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
    //  write population and turn stats in viewport.
    var txsize = p.height/40
    p.textSize(txsize)
    var vpPopString =  drawlist['cpp-pop-string']
    var cppPopString = drawlist['vp-pop-string']
    var turnString =   drawlist['turncount-string']
    var bw = ((p.textWidth(vpPopString) + p.textWidth(cppPopString) + p.textWidth(turnString)) * 2.3 ) / 3
    var bh = txsize * 6
    var bx = p.width*0.02
    var by = p.height - (p.height * 0.2)
    var br = by * 0.015
    p.push()
      p.translate(bx, by)
      p.fill(0,0,0,100)
      p.rect(0, 0, bw, bh, br)
      p.fill(255)
      p.text(vpPopString + "\n" + cppPopString + "\n" + turnString, txsize*2, txsize*2)
    p.pop()
  }

  p.windowResized = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.625
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
      drawlist['cpp-pop-string'] = rawmsg.data['cpp-pop-string']
      drawlist['vp-pop-string'] = rawmsg.data['vp-pop-string']
      drawlist['turncount-string'] = rawmsg.data['turncount-string']
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
      ['abm-rng-seedval']: parseInt($('#abm-rng-seedval').val()),
      ['abm-rng-fuzziness']: 0.1,
      ['abm-visualise-flag']: true,
      ['abm-logging-flag']: true,
      ['abm-log-frequency']: 0, // log every turn!
      ['abm-use-custom-log-filepath']: false,
      ['abm-custom-log-filepath']: "",
      ['abm-log-filepath']: "",
      ['abm-limit-duration']: false,
      ['abm-fixed-duration']: 100,
      ['abm-session-identifier']: "DebugSession"
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
