// ./build_installer.js

// 1. Import Modules
const { MSICreator } = require('electron-wix-msi');
const path = require('path');

// 2. Define input and output directory.
// Important: the directories must be absolute, not relative e.g
const APP_DIR = path.resolve(__dirname, './out/atest-desktop-win32-x64');
const OUT_DIR = path.resolve(__dirname, './windows_installer');

// 3. Instantiate the MSICreator
const msiCreator = new MSICreator({
  appDirectory: APP_DIR,
  outputDirectory: OUT_DIR,

  // Configure metadata
  description: 'API Testing Desktop Application',
  exe: 'atest',
  name: 'API Testing Desktop Application',
  manufacturer: 'LinuxSuRen',
  version: '0.0.1',

  // Configure installer User Interface
  ui: {
    chooseDirectory: true
  },
});

// 4. Create a .wxs template file
msiCreator.create().then(function(){

  // Step 5: Compile the template to a .msi file
  msiCreator.compile();
});
