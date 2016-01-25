var websocket = new WebSocket('ws://localhost:8080/ws')
var viz = new p5(sketch, 'abm-viewport')

window.onbeforeunload = function() {
  websocket.onclose = function () {}
  websocket.close()
}

websocket.onopen = function () {
  console.log('websocket connection init')
  websocket.onlose = function (e) {

  }
}

websocket.onmessage = function (e) {
  var rawmsg = JSON.parse(e.data)
  console.log(rawmsg)
  if (rawmsg.type === 'render') {
    drawlist.cpp = rawmsg.data.cpp
    drawlist.vp = rawmsg.data.vp
    drawlist.bg = rawmsg.data.bg
  }
  viz.redraw()
}

websocket.onerror = function (e) {
  console.log('ERROR: ' + e)
}

websocket.onclose = function(close) {

}
