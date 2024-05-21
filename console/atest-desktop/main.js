// main.js

// Modules to control application life and create native browser window
const { app, shell, BrowserWindow, Menu, MenuItem, ipcMain } = require('electron')
const log = require('electron-log/main');
const path = require('node:path')
const fs = require('node:fs')
const server = require('./api')
const spawn = require("child_process").spawn;
const atestHome = server.getHomeDir()

// setup log output
log.initialize();
log.transports.file.level = getLogLevel()
log.transports.file.resolvePathFn = () => server.getLogfile()

app.dock.setIcon(path.join(__dirname, "api-testing.png"))
const createWindow = () => {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 1024,
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: true,
      contextIsolation: true,
      enableRemoteModule: true
    },
    icon: path.join(__dirname, '/api-testing.ico'),
  })

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
  startServer()
  createWindow()

  app.on('activate', () => {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })

  ipcMain.on('openLogDir', () => {
    shell.openExternal('file://' + server.getLogfile())
  })
  ipcMain.on('startServer', startServer)
  ipcMain.on('stopServer', stopServer)
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
  
  if (!fs.existsSync(atestFromHome)) {
    log.info('cannot find from %s', atestFromHome)

    const data = fs.readFileSync(atestFromPkg)
    log.info('start to write file with length %d', data.length)
    
    try { 
      fs.writeFileSync(atestFromHome, data);
    } 
    catch (e) { 
      log.error('Error Code: %s', e.code); 
    }
  }
  fs.chmodSync(atestFromHome, 0o755); 

  serverProcess = spawn(atestFromHome, [
    "server",
    "--http-port", server.getPort(),
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