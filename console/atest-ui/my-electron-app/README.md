```shell
npx electron-forge import
```

```shell
npm config set registry https://registry.npmmirror.com
```

## Package

```shell
npm run package -- --platform=darwin
npm run package -- --platform=win32
npm run package -- --platform=linux
```

## For Linux

You need to install tools if you want to package Windows on Linux:
```shell
apt install wine64 zip -y
```

## Publish

export GITHUB_TOKEN=your-token

```shell
npm run publish -- --platform=darwin
npm run publish -- --platform=linux
```
