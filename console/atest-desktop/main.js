/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Modules to control application life and create native browser window
const { app, shell, BrowserWindow, Menu, MenuItem, ipcMain, contextBridge } = require('electron')
const log = require('electron-log/main');
const path = require('node:path')
const fs = require('node:fs')
const server = require('./api')
const spawn = require("child_process").spawn;
const atestHome = server.getHomeDir()
const storage = require('electron-json-storage')

// setup log output
log.initialize();

log.transports.file.level = getLogLevel()
log.transports.file.resolvePathFn = () => server.getLogfile()
if (process.platform === 'darwin'){
	app.dock.setIcon(path.join(__dirname, "api-testing.png"))
}

const windowOptions = {
  width: 1024,
  height: 600,
  frame: true,
  webPreferences: {
    preload: path.join(__dirname, 'preload.js'),
    nodeIntegration: true,
    contextIsolation: true,
    enableRemoteModule: true
  },
  icon: path.join(__dirname, '/api-testing.ico'),
}

const createWindow = () => {
  var width = storage.getSync('window.width')
  if (!isNaN(width)) {
    windowOptions.width = width
  }
  var height = storage.getSync('window.height')
  if (!isNaN(height)) {
    windowOptions.height = height
  }

  // Create the browser window.
  const mainWindow = new BrowserWindow(windowOptions)

  if (!isNaN(serverProcess.pid)) {
    // server process started by app
    mainWindow.loadURL(server.getHomePage())
  } else {
    server.control(() => {
      mainWindow.loadURL(server.getHomePage())
    }, () => {
      // and load the index.html of the app.
      mainWindow.loadFile('index.html')
    })
  }

  mainWindow.on('resize', () => {
    const size = mainWindow.getSize();
    storage.set('window.width', size[0])
    storage.set('window.height', size[1])
  })
}

const menu = new Menu()
menu.append(new MenuItem({
  label: 'Window',
  submenu: [{
    label: 'Console',
    accelerator: process.platform === 'darwin' ? 'Alt+Cmd+C' : 'Alt+Shift+C',
    click: () => {
      BrowserWindow.getFocusedWindow().loadFile('index.html');
    }
  }, {
    label: 'Server',
    accelerator: process.platform === 'darwin' ? 'Alt+Cmd+S' : 'Alt+Shift+S',
    click: () => {
      BrowserWindow.getFocusedWindow().loadURL(server.getHomePage());
    }
  }, {
    label: 'Reload',
    accelerator: process.platform === 'darwin' ? 'Cmd+R' : 'F5',
    click: () => {
      BrowserWindow.getFocusedWindow().reload()
    }
  }, {
    label: 'Developer Mode',
    accelerator: process.platform === 'darwin' ? 'Alt+Cmd+D' : 'F12',
    click: () => {
      BrowserWindow.getFocusedWindow().webContents.openDevTools();
    }
  }, {
    label: 'Quit',
    accelerator: process.platform === 'darwin' ? 'Cmd+Q' : 'Alt+Shift+Q',
    click: () => {
      app.quit()
    }
  }]
}))

Menu.setApplicationMenu(menu)

let serverProcess;
// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  ipcMain.on('openLogDir', () => {
    shell.openExternal('file://' + server.getLogfile())
  })
  ipcMain.on('startServer', startServer)
  ipcMain.on('stopServer', stopServer)
  ipcMain.on('control', (e, okCallback, errCallback) => {
    server.control(okCallback, errCallback)
  })
  ipcMain.handle('getHomePage', server.getHomePage)
  ipcMain.handle('getPort', () => {
    return server.getPort()
  })
  ipcMain.handle('getHealthzUrl', server.getHealthzUrl)

  startServer()
  createWindow()

  app.on('activate', () => {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

const startServer = () => {
  const homeData = path.join(atestHome, 'data')
  const homeBin = path.join(atestHome, 'bin')

  fs.mkdirSync(homeData, {
    recursive: true
  })
  fs.mkdirSync(homeBin, {
    recursive: true
  })

  // try to find the atest file first
  const serverFile = process.platform === "win32" ? "atest.exe" : "atest"
  const atestFromHome = path.join(homeBin, serverFile)
  const atestFromPkg = path.join(__dirname, serverFile)

  const data = fs.readFileSync(atestFromPkg)
  log.info('start to write file with length', data.length)
  
  try {
    if (process.platform === "win32") {
      const file = fs.openSync(atestFromHome, 'w');
      fs.writeSync(file, data, 0, data.length, 0);
      fs.closeSync(file);
    }else{
      fs.writeFileSync(atestFromHome, data);
    }
  } catch (e) { 
    log.error('Error Code:', e.code); 
  }
  fs.chmodSync(atestFromHome, 0o755); 

  serverProcess = spawn(atestFromHome, [
    "server",
    "--http-port", server.getPort(),
    "--port=0",
    "--local-storage", path.join(homeData, "*.yaml")
  ])
  serverProcess.stdout.on('data', (data) => {
    log.info(data.toString())
    if (data.toString().indexOf('Server is running') != -1) {
      BrowserWindow.getFocusedWindow().loadURL(server.getHomePage())
    }
  })
  serverProcess.stderr.on('data', (data) => {
    log.error(data.toString())
  })
  serverProcess.on('close', (code) => {
    log.log(`child process exited with code ${code}`);
  })
  log.info('start atest server as pid:', serverProcess.pid)
  log.info(serverProcess.spawnargs)
}

const stopServer = () => {
  if (serverProcess) {
    serverProcess.kill()
  }
}

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()

    stopServer()
  }
})
app.on('before-quit', stopServer)

function getLogLevel() {
  return 'info'
}

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.