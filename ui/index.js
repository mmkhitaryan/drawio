let socket = new WebSocket("ws://localhost:8080/ws");
var canvas = document.querySelector("#can");
var context = canvas.getContext("2d");
var connection;
canvas.width  = window.innerWidth-20;
canvas.height = window.innerHeight-20;

function drawPoint(inputx,inputy){
    console.log(inputx, inputy)
    context.fillRect(inputx, inputy, 5, 5);
}

function YourNewPoint(x,y){
    socket.send(JSON.stringify({x:x,y:y}))
    drawPoint(x,y)
}

canvas.addEventListener('mousemove', function(event) { 
    var x = event.pageX - canvas.offsetLeft;
    var y = event.pageY - canvas.offsetTop;
        
    if (event.buttons) { // If mouse movement with mouse pressed
        drawPoint(x,y)
        YourNewPoint(x,y)
    }
}, false);

socket.onmessage = function(event) {
    let data = JSON.parse(event.data)
    drawPoint(data.x,data.y)
};