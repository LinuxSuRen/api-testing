let server = require("./api.js")

server.control(() => {
    const actionBut = document.getElementById('action');
    actionBut.innerHTML = 'Stop';
})
