const { FusesPlugin } = require('@electron-forge/plugin-fuses');
const { FuseV1Options, FuseVersion } = require('@electron/fuses');

module.exports = {
  packagerConfig: {
    icon: 'api-testing.ico',
    asar: true,
  },
  rebuildConfig: {},
  makers: [
    {
      name: '@electron-forge/maker-squirrel',
      config: {
        icon: 'api-testing.ico'
      },
    },
    {
      name: '@electron-forge/maker-deb',
      config: {
        options: {
          icon: 'api-testing.ico'
        }
      },
    },
    {
      name: '@electron-forge/maker-rpm',
      config: {
        icon: 'api-testing.ico'
      },
    },
    {
      name: '@electron-forge/maker-dmg',
      config: {
        format: 'ULFO',
        icon: 'api-testing.ico'
      }
    },
    {
      name: '@electron-forge/maker-wix',
      config: {
        language: 1033,
        manufacturer: 'API Testing Authors',
        icon: 'api-testing.ico'
      },
      ui: {
        "enabled": true,
        "chooseDirectory": true
      }
    }
  ],
  plugins: [
    {
      name: '@electron-forge/plugin-auto-unpack-natives',
      config: {},
    },
    // Fuses are used to enable/disable various Electron functionality
    // at package time, before code signing the application
    new FusesPlugin({
      version: FuseVersion.V1,
      [FuseV1Options.RunAsNode]: false,
      [FuseV1Options.EnableCookieEncryption]: true,
      [FuseV1Options.EnableNodeOptionsEnvironmentVariable]: false,
      [FuseV1Options.EnableNodeCliInspectArguments]: false,
      [FuseV1Options.EnableEmbeddedAsarIntegrityValidation]: true,
      [FuseV1Options.OnlyLoadAppFromAsar]: true,
    }),
  ],
  publishers: [
    {
      name: '@electron-forge/publisher-github',
      config: {
        repository: {
          owner: 'linuxsuren',
          name: 'api-testing'
        },
        prerelease: true,
        force: true
      }
    }
  ]
};
