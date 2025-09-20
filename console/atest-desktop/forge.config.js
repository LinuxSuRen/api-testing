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

const { FusesPlugin } = require('@electron-forge/plugin-fuses');
const { FuseV1Options, FuseVersion } = require('@electron/fuses');

module.exports = {
  packagerConfig: {
    icon: './assets/icons/atest',
    asar: true
  },
  rebuildConfig: {},
  makers: [
    {
      name: '@electron-forge/maker-squirrel',
      icon: './assets/icons/atest.ico'
    },
    {
      name: '@electron-forge/maker-deb',
      icon: './assets/icons/atest.png'
    },
    {
      name: '@electron-forge/maker-rpm',
      icon: './assets/icons/atest.png'
    },
    {
      name: '@electron-forge/maker-dmg',
      icon: './assets/icons/atest.icns',
      background: './assets/icons/atest.png',
      format: 'ULFO'
    },
    {
      name: '@electron-forge/maker-wix',
      config: {
        language: 1033,
        manufacturer: 'API Testing Authors',
        icon: './assets/icons/atest.ico',
        ui: {
          "enabled": true,
          "chooseDirectory": true
        },
        beforeCreate: (msiCreator) => {
          // Add installation directory to system PATH
          msiCreator.wixTemplate = msiCreator.wixTemplate.replace(
            '</Product>',
            `    <Property Id="ARPINSTALLLOCATION" Value="[INSTALLDIR]" />
    <CustomAction Id="AddToPath" Property="PATH" Value="[INSTALLDIR]" Execute="immediate" />
    <CustomAction Id="RemoveFromPath" Property="PATH" Value="[INSTALLDIR]" Execute="immediate" />

    <InstallExecuteSequence>
      <Custom Action="AddToPath" After="InstallFiles">NOT Installed</Custom>
      <Custom Action="RemoveFromPath" Before="RemoveFiles">REMOVE="ALL"</Custom>
    </InstallExecuteSequence>

    <Component Id="PathComponent" Guid="*" Directory="INSTALLDIR">
      <Environment Id="PATH" Name="PATH" Value="[INSTALLDIR]" Permanent="no" Part="last" Action="set" System="yes" />
    </Component>

    </Product>`
          );

          // Ensure INSTALLDIR is properly defined in the directory structure
          msiCreator.wixTemplate = msiCreator.wixTemplate.replace(
            '<Directory Id="TARGETDIR" Name="SourceDir">',
            `<Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="ProgramFilesFolder">
        <Directory Id="INSTALLDIR" Name="API Testing" />
      </Directory>`
          );
        }
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
