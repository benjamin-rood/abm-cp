
var drawlist = new DrawList(initDrawObj)
var wsUrl = "ws://" + window.location.hostname + ":" + window.location.port + "/ws"
var vizSocket = new WebSocket(wsUrl)
var viz = new p5(sketch, 'abm-viewport')

window.onbeforeunload = function() {
  vizSocket.close()
}

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
