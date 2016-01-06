var ws = new WebSocket("ws://localhost:8080/ws")

wsocket.onopen = function() {
  console.log("websocket connection init")
}

wsocket.onmessage = function(e) {
  var rawmsg = e.data
  if (rawmsg.type == "render") {
    renderlist = rawmsg.data
  }
  console.log("websocket message received")

};

wsocket.onerror = function(e) {
  console.log("ERROR: " + e)
}

function sendData() {
  // var insObj = {
  //   'seedVal': seed,
  //   'rangeVal': range
  // };
  // var msg = JSON.stringify(insObj);
  // wsocket.send(msg);
}
