var wsocket = new WebSocket('ws://localhost:8080/ws')

wsocket.onopen = function () {
  console.log('websocket connection init')
  var viz = new p5(sketch, 'abm-viewport')
}

wsocket.onmessage = function (e) {
  var rawmsg = JSON.parse(e.data)
  // console.log(rawmsg)
  if (rawmsg.type === 'render') {
    drawlist = new DrawList(rawmsg)
  }
  // if (rawmsg.type === 'statistics') {
  //
  // }
}

wsocket.onerror = function (e) {
  console.log('ERROR: ' + e)
}
