const exp = require('constants')
const path = require('node:path')

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

function getHomeDir() {
    const homedir = require('os').homedir();
    return path.join(homedir, ".config", 'atest')
}

function getLogfile() {
    return path.join(getHomeDir(), 'log.log')
}

exports.getPort = getPort
exports.getHomePage = getHomePage
exports.getHomeDir = getHomeDir
exports.getLogfile = getLogfile
