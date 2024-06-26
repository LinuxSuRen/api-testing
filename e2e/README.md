You can build the image locally in the repository root directory:

If you are a Linux/MacOS user, you can use the following command:
```shell
REGISTRY=ghcr.io TAG=master make image
```

If you are a Windows user, you can use the following command:
```powershell
Set-Content -Path "env:REGISTRY" -Value "ghcr.io"
Set-Content -Path "env:TAG" -Value "master"
make image
```
