
var conn = new WebSocket("ws://localhost:3000/ws-chat")

conn.onmessage = function(m) {
    console.log("recieved: " + m.data)
}

function sendMessage() {
  var msg = document.getElementById("msg-input").value
  conn.send(msg)
}
