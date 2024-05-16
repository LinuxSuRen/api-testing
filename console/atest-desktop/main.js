// main.js

// Modules to control application life and create native browser window
const { app, BrowserWindow, Menu, MenuItem } = require('electron')
const path = require('node:path')
const server = require('./api')
const spawn = require("child_process").spawn;

const createWindow = () => {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 1024,
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: true,
      contextIsolation: false,
      enableRemoteModule: true
    },
    icon: path.join(__dirname, '/api-testing.ico'),
  })

  server.control(() => {
    mainWindow.loadURL(server.getHomePage())
  }, () => {
    // and load the index.html of the app.
    mainWindow.loadFile('index.html')
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
  const homedir = require('os').homedir();

  serverProcess = spawn("atest", [
    "server",
    "--http-port", server.getPort(),
    // TODO below setting is not working
    "--local-storage", path.join(homedir, ".atest", "data", "*.yaml")
  ]);

  createWindow()

  app.on('activate', () => {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()

    if (serverProcess) {
      serverProcess.kill();
    }
  }
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.