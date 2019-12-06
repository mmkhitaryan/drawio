let socket = new WebSocket("ws://localhost:8080/ws");
var canvas = document.querySelector("#can");
var context = canvas.getContext("2d");

let isDrawing = false;
let x = 0;
let y = 0;

canvas.width  = window.innerWidth-20;
canvas.height = window.innerHeight-20;

function drawPoint(new_x, new_y){
    old_x = x
    old_y = y

    context.beginPath();
    context.strokeStyle = 'black';
    context.lineWidth = 5;
    context.lineCap = "round";
    context.moveTo(new_x, new_y);
    context.lineTo(old_x, old_y);
    context.stroke();
    context.closePath();
}

function YourNewPoint(x,y){
    socket.send(JSON.stringify({
        new_x: x,
        new_y: y,
        old_x: x,
        old_y: y
    })) // add also new and old vars 
    drawPoint(x,y)
}

canvas.addEventListener('mousedown', function(event) { 
    x = event.pageX - canvas.offsetLeft;
    y = event.pageY - canvas.offsetTop;
    isDrawing = true;
});

canvas.addEventListener('mousemove', function(event) {
    if (isDrawing) {
        new_x = event.pageX - canvas.offsetLeft;
        new_y = event.pageY - canvas.offsetTop;
        YourNewPoint(new_x, new_y)
        x = new_x
        y = new_y
    }
}, false);

canvas.addEventListener('mouseup', function(event) {
    if (isDrawing === true) {
        x = 0;
        y = 0;
        isDrawing = false;
    }
})

//socket.onmessage = function(event) {
//   let data = JSON.parse(event.data)
//    drawPoint(data.x,data.y)
//};