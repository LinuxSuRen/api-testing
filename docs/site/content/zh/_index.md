---
title: API Testing
---

{{< blocks/cover title="欢迎访问 API Testing！" image_anchor="top" height="full" >}}
<a class="btn btn-lg btn-primary me-3 mb-4" href="{{< relref "/latest" >}}">
  开始使用 <i class="fas fa-arrow-alt-circle-right ms-2"></i>
</a>
<a class="btn btn-lg btn-secondary me-3 mb-4" href="{{< relref "/contributions" >}}">
  参与贡献 <i class="fa fa-heartbeat ms-2 "></i>
</a>
<p class="lead mt-5">开源接口调试 & 测试工具</p>
<!-- 向下翻页图标 -->
{{< blocks/link-down color="white" >}}
<p class="lead mt-5">让我们一起提高代码质量吧！</p>
{{< /blocks/cover >}}

{{< blocks/section >}}

<div id="download-links" style="text-align: center;">
  <div class="mb-4">
    下载桌面应用
  </div>

  <a class="btn btn-primary me-2 mb-2" id="download-windows" href="https://github.com/linuxsuren/api-testing/releases/latest/download/atest-desktop.msi" style="display:none;">
    Windows
  </a>
  <a class="btn btn-secondary me-2 mb-2" id="download-windows-proxy" href="https://files.m.daocloud.io/github.com/LinuxSuRen/api-testing/releases/download/v0.0.19/atest-desktop.msi" style="display:none;">
    DaoCloud 镜像站
  </a>
  <a class="btn btn-primary me-2 mb-2" id="download-macos" href="https://github.com/linuxsuren/api-testing/releases/latest/download/atest-desktop-0.0.19-x64.dmg" style="display:none;">
    macOS
  </a>
  <a class="btn btn-secondary me-2 mb-2" id="download-macos-proxy" href="https://files.m.daocloud.io/github.com/LinuxSuRen/api-testing/releases/download/v0.0.19/atest-desktop-0.0.19-x64.dmg" style="display:none;">
    DaoCloud 镜像站
  </a>
  <a class="btn btn-primary me-2 mb-2" id="download-linux" href="https://github.com/linuxsuren/api-testing/releases/latest/download/atest-desktop_0.0.19_amd64.deb" style="display:none;">
    Linux
  </a>
  <a class="btn btn-secondary me-2 mb-2" id="download-linux-proxy" href="https://files.m.daocloud.io/github.com/LinuxSuRen/api-testing/releases/download/v0.0.19/atest-desktop_0.0.19_amd64.deb" style="display:none;">
    DaoCloud 镜像站
  </a>
</div>
<noscript>
  <p>请根据您的操作系统选择以下链接手动下载：</p>
  <ul>
    <li>
      <a href="https://github.com/linuxsuren/api-testing/releases/latest/download/atest-desktop.msi">Windows</a> |
      <a href="https://files.m.daocloud.io/github.com/LinuxSuRen/api-testing/releases/download/v0.0.19/atest-desktop.msi">代理下载</a>
    </li>
    <li>
      <a href="https://github.com/linuxsuren/api-testing/releases/latest/download/atest-desktop-0.0.19-x64.dmg">macOS</a> |
      <a href="https://files.m.daocloud.io/github.com/LinuxSuRen/api-testing/releases/download/v0.0.19/atest-desktop-0.0.19-x64.dmg">代理下载</a>
    </li>
    <li>
      <a href="https://github.com/linuxsuren/api-testing/releases/latest/download/atest-desktop_0.0.19_amd64.deb">Linux</a> |
      <a href="https://files.m.daocloud.io/github.com/LinuxSuRen/api-testing/releases/download/v0.0.19/atest-desktop_0.0.19_amd64.deb">代理下载</a>
    </li>
  </ul>
</noscript>
<script>
  (function() {
    var platform = window.navigator.platform.toLowerCase();
    if (platform.indexOf('win') >= 0) {
      document.getElementById('download-windows').style.display = 'inline-block';
      document.getElementById('download-windows-proxy').style.display = 'inline-block';
    } else if (platform.indexOf('mac') >= 0) {
      document.getElementById('download-macos').style.display = 'inline-block';
      document.getElementById('download-macos-proxy').style.display = 'inline-block';
    } else if (platform.indexOf('linux') >= 0) {
      document.getElementById('download-linux').style.display = 'inline-block';
      document.getElementById('download-linux-proxy').style.display = 'inline-block';
    } else {
      // 默认全部显示
      document.getElementById('download-windows').style.display = 'inline-block';
      document.getElementById('download-windows-proxy').style.display = 'inline-block';
      document.getElementById('download-macos').style.display = 'inline-block';
      document.getElementById('download-macos-proxy').style.display = 'inline-block';
      document.getElementById('download-linux').style.display = 'inline-block';
      document.getElementById('download-linux-proxy').style.display = 'inline-block';
    }
  })();
</script>
{{< /blocks/section >}}