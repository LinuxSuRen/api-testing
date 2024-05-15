exports.control = function(okCallback, errorCallback) {
    fetch('http://localhost:' + getPort() + '/healthz').
        then(okCallback).catch(errorCallback)
}

function getPort() {
    // TODO support set this value
    return 7788
}

function getHomePage() {
    return 'http://localhost:' + getPort()
}

exports.getPort = getPort
exports.getHomePage = getHomePage
