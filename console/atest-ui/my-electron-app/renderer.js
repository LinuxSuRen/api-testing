let spawn = require("child_process").spawn;
let server = require("./api.js")

server.control(() => {
    const actionBut = document.getElementById('action');
    actionBut.innerHTML = 'Stop';
})

let process;
function start() {
    process = spawn("atest", [
        "server",
        "--http-port",
        server.getPort()
    ]);
    
    process.stdout.on("data", (data) => {
        console.log(data.toString());
    });
    
    process.stderr.on("data", (err) => {
        console.log(err.toString());
    });
    
    process.on("exit", (code) => {
        console.log(code);
    });
    load();
}

function stop() {
    if (process) {
        process.kill();
    }
    load();
}
