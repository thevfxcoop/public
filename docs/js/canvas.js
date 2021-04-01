
const gridSize = 16;
const someLetters = "ABDEGHIJKLMNQRSTUWYZ";
const otherLetters = "VFXCOOP";

// Generate random letter and color
function letter() {
    var randLetter = Math.round(Math.random() * (someLetters.length + otherLetters.length - 1));
    if(randLetter < otherLetters.length) {
        return [ otherLetters[randLetter], "#888" ];
    } else {
        return [ someLetters[randLetter - otherLetters.length], "#555" ];
    }
}

function drawCanvas(canvas) {
    const ctx = canvas.getContext('2d');
    const center = gridSize >> 1;

    // Calculate offset
    const xoffset = (canvas.width % gridSize) >> 2 + 0.5;
    const yoffset = (canvas.height % gridSize) >> 2 + 0.5;
    ctx.translate(xoffset,yoffset);

    // Set styles
    ctx.font = '' + center + 'px "Lato",sans-serif';
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillStyle = "#555";

    // Draw grid
    for (var y = 0; y < canvas.clientHeight; y += gridSize) {
        for (var x = 0; x < canvas.clientWidth; x += gridSize) {
            var text = letter();
            ctx.fillStyle = text[1];
            ctx.fillText(text[0], x - center, y - center);
        }
    }
}

window.addEventListener("DOMContentLoaded", (event) => {
    var canvas = document.querySelector("#canvas");
    canvas.width = canvas.clientWidth;
    canvas.height = canvas.clientHeight;
    drawCanvas(document.querySelector("#canvas"));
});

window.addEventListener('resize',(event) => {
    var canvas = document.querySelector("#canvas");
    canvas.width = canvas.clientWidth;
    canvas.height = canvas.clientHeight;
    drawCanvas(document.querySelector("#canvas"));
});
