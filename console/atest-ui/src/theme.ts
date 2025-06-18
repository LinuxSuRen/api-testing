const defaultTheme = {
  "color": {
    "--el-color-white": "#ffffff",
      "--el-color-black": "#000000",
      "--el-color-primary": "#409eff",
      "--el-color-primary-light-3": "#79bbff",
      "--el-color-primary-light-5": "#a0cfff",
      "--el-color-primary-light-7": "#c6e2ff",
      "--el-color-primary-light-8": "#d9ecff",
      "--el-color-primary-light-9": "#ecf5ff",
      "--el-color-primary-dark-2": "#337ecc",
      "--el-color-success": "#67c23a",
      "--el-color-success-light-3": "#95d475",
      "--el-color-success-light-5": "#b3e19d",
      "--el-color-success-light-7": "#d1edc4",
      "--el-color-success-light-8": "#e1f3d8",
      "--el-color-success-light-9": "#f0f9eb",
      "--el-color-success-dark-2": "#529b2e",
      "--el-color-warning": "#e6a23c",
      "--el-color-warning-light-3": "#eebe77",
      "--el-color-warning-light-5": "#f3d19e",
      "--el-color-warning-light-7": "#f8e3c5",
      "--el-color-warning-light-8": "#faecd8",
      "--el-color-warning-light-9": "#fdf6ec",
      "--el-color-warning-dark-2": "#b88230",
      "--el-color-danger": "#f56c6c",
      "--el-color-danger-light-3": "#f89898",
      "--el-color-danger-light-5": "#fab6b6",
      "--el-color-danger-light-7": "#fcd3d3",
      "--el-color-danger-light-8": "#fde2e2",
      "--el-color-danger-light-9": "#fef0f0",
      "--el-color-danger-dark-2": "#c45656",
      "--el-color-error": "#f56c6c",
      "--el-color-info": "#909399",
      "--el-color-info-light-3": "#b1b3b8",
      "--el-color-info-light-5": "#c8c9cc",
      "--el-color-info-light-7": "#dedfe0",
      "--el-color-info-light-8": "#e9e9eb",
      "--el-color-info-light-9": "#f4f4f5",
      "--el-bg-color": "#ffffff",
      "--el-bg-color-page": "#f2f3f5",
      "--el-bg-color-overlay": "#ffffff",
      "--el-text-color-primary": "#303133",
      "--el-text-color-regular": "#606266",
      "--el-text-color-secondary": "#909399",
      "--el-text-color-placeholder": "#a8abb2",
      "--el-text-color-disabled": "#c0c4cc",
      "--el-border-color": "#dcdfe6",
      "--el-border-color-light": "#e4e7ed",
      "--el-border-color-lighter": "#ebeef5",
      "--el-border-color-extra-light": "#f2f6fc",
      "--el-border-color-dark": "#d4d7de",
      "--el-border-color-darker": "#cdd0d6",
      "--el-fill-color": "#f0f2f5",
      "--el-fill-color-light": "#f5f7fa",
      "--el-fill-color-lighter": "#fafafa",
      "--el-fill-color-extra-light": "#fafcff",
      "--el-fill-color-dark": "#ebedf0",
      "--el-fill-color-darker": "#e6e8eb",
      "--el-fill-color-blank": "#ffffff"
  },
  "common": {
    "--el-border-width": "1px",
      "--el-border-style": "solid",
      "--el-border-color-hover": "",
      "--el-border": "var(--el-border-width) var(--el-border-style) var(--el-border-color)",
      "--el-svg-monochrome-grey": "#dcdde0",
      "--el-border-radius-base": "4px",
      "--el-border-radius-small": "2px",
      "--el-border-radius-round": "20px",
      "--el-border-radius-circle": "100%",
      "--el-box-shadow": "0px 12px 32px 4px rgba(0, 0, 0, 0.04), 0px 8px 20px rgba(0, 0, 0, 0.08)",
      "--el-box-shadow-light": "0px 0px 12px rgba(0, 0, 0, 0.12)",
      "--el-box-shadow-lighter": "0px 0px 6px rgba(0, 0, 0, 0.12)",
      "--el-box-shadow-dark": "0px 16px 48px 16px rgba(0, 0, 0, 0.08), 0px 12px 32px rgba(0, 0, 0, 0.12), 0px 8px 16px -8px rgba(0, 0, 0, 0.16)"
  },
  "font": {
    "--el-font-size-extra-large": "20px",
      "--el-font-size-large": "18px",
      "--el-font-size-medium": "16px",
      "--el-font-size-base": "14px",
      "--el-font-size-small": "13px",
      "--el-font-size-extra-small": "12px",
      "--el-font-family": "'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif",
      "--el-font-weight-primary": "500",
      "--el-font-line-height-primary": "24px"
  },
  "size": {
    "--el-component-size-large": "40px",
      "--el-component-size": "32px",
      "--el-component-size-small": "24px"
  },
  "z-index": {
    "--el-index-normal": "1",
      "--el-index-top": "1000",
      "--el-index-popper": "2000"
  },
  "components": {
    "button": {
      "--el-button-font-weight": "var(--el-font-weight-primary)",
        "--el-button-border-color": "var(--el-border-color)",
        "--el-button-bg-color": "var(--el-fill-color-blank)",
        "--el-button-text-color": "var(--el-text-color-regular)",
        "--el-button-disabled-text-color": "var(--el-disabled-text-color)",
        "--el-button-disabled-bg-color": "var(--el-fill-color-blank)",
        "--el-button-disabled-border-color": "var(--el-border-color-light)",
        "--el-button-hover-text-color": "var(--el-color-primary)",
        "--el-button-hover-bg-color": "var(--el-color-primary-light-9)",
        "--el-button-hover-border-color": "var(--el-color-primary-light-7)",
        "--el-button-active-text-color": "var(--el-button-hover-text-color)",
        "--el-button-active-border-color": "var(--el-color-primary)",
        "--el-button-active-bg-color": "var(--el-button-hover-bg-color)",
        "--el-button-outline-color": "var(--el-color-primary-light-5)"
    },
    "input": {
      "--el-input-text-color": "var(--el-text-color-regular)",
        "--el-input-border": "var(--el-border)",
        "--el-input-hover-border": "var(--el-border-color-hover)",
        "--el-input-focus-border": "var(--el-color-primary)",
        "--el-input-transparent-border": "0 0 0 1px transparent inset",
        "--el-input-border-color": "var(--el-border-color)",
        "--el-input-border-radius": "var(--el-border-radius-base)",
        "--el-input-bg-color": "var(--el-fill-color-blank)",
        "--el-input-icon-color": "var(--el-text-color-placeholder)",
        "--el-input-placeholder-color": "var(--el-text-color-placeholder)",
        "--el-input-hover-border-color": "var(--el-border-color-hover)",
        "--el-input-clear-hover-color": "var(--el-text-color-secondary)",
        "--el-input-focus-border-color": "var(--el-color-primary)"
    },
    "tag": {
      "--el-tag-bg-color": "var(--el-color-primary-light-9)",
        "--el-tag-text-color": "var(--el-color-primary)",
        "--el-tag-border-color": "var(--el-color-primary-light-8)",
        "--el-tag-hover-color": "var(--el-color-primary)"
    },
    "notification": {
      "--el-notification-width": "330px",
        "--el-notification-padding": "14px 26px 14px 13px",
        "--el-notification-radius": "8px",
        "--el-notification-shadow": "var(--el-box-shadow-light)",
        "--el-notification-border-color": "var(--el-border-color-lighter)",
        "--el-notification-icon-size": "24px",
        "--el-notification-close-font-size": "16px",
        "--el-notification-group-margin-left": "13px",
        "--el-notification-group-margin-right": "8px",
        "--el-notification-content-font-size": "13px",
        "--el-notification-content-color": "var(--el-text-color-regular)",
        "--el-notification-title-font-size": "16px",
        "--el-notification-title-color": "var(--el-text-color-primary)",
        "--el-notification-close-color": "var(--el-text-color-secondary)",
        "--el-notification-close-hover-color": "var(--el-text-color-regular)"
    }
  }
}

const bubblegumTheme = {
  "color": {
    "--el-color-white": "#ffffff",
    "--el-color-black": "#000000",
    // 主色调：泡泡糖粉 (原蓝色系改为粉紫色系)
    "--el-color-primary": "#FF6BBD", // 亮粉色
    "--el-color-primary-light-3": "#FF8FCC",
    "--el-color-primary-light-5": "#FFB3DB",
    "--el-color-primary-light-7": "#FFD7EA",
    "--el-color-primary-light-8": "#FFE9F4",
    "--el-color-primary-light-9": "#FFF5FA",
    "--el-color-primary-dark-2": "#E05AA8",
    // 成功色：糖果绿
    "--el-color-success": "#5CD1C3", // 薄荷绿
    "--el-color-success-light-3": "#85DECE",
    "--el-color-success-light-5": "#AEEADF",
    "--el-color-success-light-7": "#D6F6EF",
    "--el-color-success-light-8": "#E7FAF5",
    "--el-color-success-light-9": "#F3FDFB",
    "--el-color-success-dark-2": "#4AB9AA",
    // 警告色：糖果黄
    "--el-color-warning": "#FFD166", // 奶油黄
    "--el-color-warning-light-3": "#FFDC8F",
    "--el-color-warning-light-5": "#FFE8B3",
    "--el-color-warning-light-7": "#FFF3D9",
    "--el-color-warning-light-8": "#FFF8E6",
    "--el-color-warning-light-9": "#FFFCF2",
    "--el-color-warning-dark-2": "#E6BC5C",
    // 危险色：糖果红
    "--el-color-danger": "#FF7AA2", // 草莓红
    "--el-color-danger-light-3": "#FF9FB9",
    "--el-color-danger-light-5": "#FFC0D3",
    "--el-color-danger-light-7": "#FFDFE9",
    "--el-color-danger-light-8": "#FFECF2",
    "--el-color-danger-light-9": "#FFF6F9",
    "--el-color-danger-dark-2": "#E66B92",
    // 信息色：糖果紫
    "--el-color-error": "#FF7AA2", // 同danger
    "--el-color-info": "#A98CFF", // 薰衣草紫
    "--el-color-info-light-3": "#C0ABFF",
    "--el-color-info-light-5": "#D5C7FF",
    "--el-color-info-light-7": "#EAE3FF",
    "--el-color-info-light-8": "#F2EEFF",
    "--el-color-info-light-9": "#F9F7FF",
    // 背景与文本
    "--el-bg-color": "#FFFFFF",
    "--el-bg-color-page": "#FFF5FA", // 浅粉背景
    "--el-bg-color-overlay": "#FFFFFF",
    "--el-text-color-primary": "#4A2E5C", // 深紫代替黑色
    "--el-text-color-regular": "#7A5C8D",
    "--el-text-color-secondary": "#B8A6C6",
    "--el-text-color-placeholder": "#D4CADE",
    "--el-text-color-disabled": "#E5DDEB",
    // 边框色
    "--el-border-color": "#F0D5E7",
    "--el-border-color-light": "#F5E3F0",
    "--el-border-color-lighter": "#FAF0F8",
    "--el-border-color-extra-light": "#FDF9FC",
    "--el-border-color-dark": "#E5C2DA",
    "--el-border-color-darker": "#D9B0CE",
    // 填充色
    "--el-fill-color": "#FFF0F7",
    "--el-fill-color-light": "#FFF5FA",
    "--el-fill-color-lighter": "#FFFAFD",
    "--el-fill-color-extra-light": "#FFFDFE",
    "--el-fill-color-dark": "#FFE5F2",
    "--el-fill-color-darker": "#FFD5EB",
    "--el-fill-color-blank": "#FFFFFF"
  },
  "common": {
    "--el-border-width": "1px",
    "--el-border-style": "solid",
    "--el-border-color-hover": "",
    "--el-border": "var(--el-border-width) var(--el-border-style) var(--el-border-color)",
    "--el-svg-monochrome-grey": "#dcdde0",
    "--el-border-radius-base": "4px",
    "--el-border-radius-small": "2px",
    "--el-border-radius-round": "20px",
    "--el-border-radius-circle": "100%",
    "--el-box-shadow": "0px 12px 32px 4px rgba(0, 0, 0, 0.04), 0px 8px 20px rgba(0, 0, 0, 0.08)",
    "--el-box-shadow-light": "0px 0px 12px rgba(0, 0, 0, 0.12)",
    "--el-box-shadow-lighter": "0px 0px 6px rgba(0, 0, 0, 0.12)",
    "--el-box-shadow-dark": "0px 16px 48px 16px rgba(0, 0, 0, 0.08), 0px 12px 32px rgba(0, 0, 0, 0.12), 0px 8px 16px -8px rgba(0, 0, 0, 0.16)"
  },
  "font": {
    "--el-font-size-extra-large": "20px",
    "--el-font-size-large": "18px",
    "--el-font-size-medium": "16px",
    "--el-font-size-base": "14px",
    "--el-font-size-small": "13px",
    "--el-font-size-extra-small": "12px",
    "--el-font-family": "'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif",
    "--el-font-weight-primary": "500",
    "--el-font-line-height-primary": "24px"
  },
  "size": {
    "--el-component-size-large": "40px",
    "--el-component-size": "32px",
    "--el-component-size-small": "24px"
  },
  "z-index": {
    "--el-index-normal": "1",
    "--el-index-top": "1000",
    "--el-index-popper": "2000"
  },
  "components": {
    "button": {
      "--el-button-font-weight": "var(--el-font-weight-primary)",
      "--el-button-border-color": "var(--el-border-color)",
      "--el-button-bg-color": "var(--el-fill-color-blank)",
      "--el-button-text-color": "var(--el-text-color-regular)",
      "--el-button-disabled-text-color": "var(--el-disabled-text-color)",
      "--el-button-disabled-bg-color": "var(--el-fill-color-blank)",
      "--el-button-disabled-border-color": "var(--el-border-color-light)",
      "--el-button-hover-text-color": "var(--el-color-primary)",
      "--el-button-hover-bg-color": "var(--el-color-primary-light-9)",
      "--el-button-hover-border-color": "var(--el-color-primary-light-7)",
      "--el-button-active-text-color": "var(--el-button-hover-text-color)",
      "--el-button-active-border-color": "var(--el-color-primary)",
      "--el-button-active-bg-color": "var(--el-button-hover-bg-color)",
      "--el-button-outline-color": "var(--el-color-primary-light-5)"
    },
    "input": {
      "--el-input-text-color": "var(--el-text-color-regular)",
      "--el-input-border": "var(--el-border)",
      "--el-input-hover-border": "var(--el-border-color-hover)",
      "--el-input-focus-border": "var(--el-color-primary)",
      "--el-input-transparent-border": "0 0 0 1px transparent inset",
      "--el-input-border-color": "var(--el-border-color)",
      "--el-input-border-radius": "var(--el-border-radius-base)",
      "--el-input-bg-color": "var(--el-fill-color-blank)",
      "--el-input-icon-color": "var(--el-text-color-placeholder)",
      "--el-input-placeholder-color": "var(--el-text-color-placeholder)",
      "--el-input-hover-border-color": "var(--el-border-color-hover)",
      "--el-input-clear-hover-color": "var(--el-text-color-secondary)",
      "--el-input-focus-border-color": "var(--el-color-primary)"
    },
    "tag": {
      "--el-tag-bg-color": "var(--el-color-primary-light-9)",
      "--el-tag-text-color": "var(--el-color-primary)",
      "--el-tag-border-color": "var(--el-color-primary-light-8)",
      "--el-tag-hover-color": "var(--el-color-primary)"
    },
    "notification": {
      "--el-notification-width": "330px",
      "--el-notification-padding": "14px 26px 14px 13px",
      "--el-notification-radius": "8px",
      "--el-notification-shadow": "var(--el-box-shadow-light)",
      "--el-notification-border-color": "var(--el-border-color-lighter)",
      "--el-notification-icon-size": "24px",
      "--el-notification-close-font-size": "16px",
      "--el-notification-group-margin-left": "13px",
      "--el-notification-group-margin-right": "8px",
      "--el-notification-content-font-size": "13px",
      "--el-notification-content-color": "var(--el-text-color-regular)",
      "--el-notification-title-font-size": "16px",
      "--el-notification-title-color": "var(--el-text-color-primary)",
      "--el-notification-close-color": "var(--el-text-color-secondary)",
      "--el-notification-close-hover-color": "var(--el-text-color-regular)"
    }
  }
}

const vintagePaperTheme = {
  "color": {
    "--el-color-white": "#F8F3E6",  // 米白色纸张底色
    "--el-color-black": "#3C2F2D",  // 深褐色（复古文字色）
    "--el-color-primary": "#a67c52",  // 墨绿色（复古主色）
    "--el-color-primary-light-3": "#708A4C",
    "--el-color-primary-light-5": "#8DA269",
    "--el-color-primary-light-7": "#A9BA8D",
    "--el-color-primary-light-8": "#C1CEAD",
    "--el-color-primary-light-9": "#D8E0CE",
    "--el-color-primary-dark-2": "#435426",
    "--el-color-success": "#6B8E23",  // 橄榄绿
    "--el-color-success-light-3": "#87A64D",
    "--el-color-success-light-5": "#A3BE77",
    "--el-color-success-light-7": "#BFD4A3",
    "--el-color-success-light-8": "#D3E0C1",
    "--el-color-success-light-9": "#E7EEDF",
    "--el-color-success-dark-2": "#56711C",
    "--el-color-warning": "#B68D40",  // 琥珀色（旧铜色）
    "--el-color-warning-light-3": "#C9A566",
    "--el-color-warning-light-5": "#D8BD8C",
    "--el-color-warning-light-7": "#E6D5B3",
    "--el-color-warning-light-8": "#EFE3C9",
    "--el-color-warning-light-9": "#F7F1E4",
    "--el-color-warning-dark-2": "#937133",
    "--el-color-danger": "#8B4513",  // 马鞍棕（旧皮革色）
    "--el-color-danger-light-3": "#A66D42",
    "--el-color-danger-light-5": "#C19671",
    "--el-color-danger-light-7": "#DCBFA3",
    "--el-color-danger-light-8": "#E9D4C2",
    "--el-color-danger-light-9": "#F5E9E0",
    "--el-color-danger-dark-2": "#70370F",
    "--el-color-error": "#8B4513",
    "--el-color-info": "#7D7C7A",  // 旧报纸灰
    "--el-color-info-light-3": "#9A9997",
    "--el-color-info-light-5": "#B5B4B2",
    "--el-color-info-light-7": "#D0CFCD",
    "--el-color-info-light-8": "#DFDEDD",
    "--el-color-info-light-9": "#EEEDEC",
    // 背景色系统（统一为纸张色调）
    "--el-bg-color": "#F8F3E6",  // 主背景
    "--el-bg-color-page": "#EFE9D9",  // 页面背景（稍深）
    "--el-bg-color-overlay": "#FDF9F0",  // 浮层背景（更浅）
    // 文字颜色系统（深棕/灰色系）
    "--el-text-color-primary": "#3C2F2D",  // 主要文字
    "--el-text-color-regular": "#5D4C47",
    "--el-text-color-secondary": "#7D716C",
    "--el-text-color-placeholder": "#A09892",
    "--el-text-color-disabled": "#C5BEB9",
    // 边框色系统（浅褐/米黄）
    "--el-border-color": "#D8CBBF",  // 基础边框
    "--el-border-color-light": "#E0D6CC",
    "--el-border-color-lighter": "#E8E1D9",
    "--el-border-color-extra-light": "#F0ECE5",
    "--el-border-color-dark": "#D0C0B0",
    "--el-border-color-darker": "#C8B6A6",
    // 填充色系统（纸张衍生色）
    "--el-fill-color": "#EFE9D9",
    "--el-fill-color-light": "#F4F0E3",
    "--el-fill-color-lighter": "#F9F6ED",
    "--el-fill-color-extra-light": "#FDFBF7",
    "--el-fill-color-dark": "#E6DECD",
    "--el-fill-color-darker": "#DED3BC",
    "--el-fill-color-blank": "#FDF9F0"
  },
  "common": {
    "--el-border-width": "1px",
    "--el-border-style": "solid",
    "--el-border-color-hover": "",
    "--el-border": "var(--el-border-width) var(--el-border-style) var(--el-border-color)",
    "--el-svg-monochrome-grey": "#dcdde0",
    "--el-border-radius-base": "4px",
    "--el-border-radius-small": "2px",
    "--el-border-radius-round": "20px",
    "--el-border-radius-circle": "100%",
    "--el-box-shadow": "0px 12px 32px 4px rgba(0, 0, 0, 0.04), 0px 8px 20px rgba(0, 0, 0, 0.08)",
    "--el-box-shadow-light": "0px 0px 12px rgba(0, 0, 0, 0.12)",
    "--el-box-shadow-lighter": "0px 0px 6px rgba(0, 0, 0, 0.12)",
    "--el-box-shadow-dark": "0px 16px 48px 16px rgba(0, 0, 0, 0.08), 0px 12px 32px rgba(0, 0, 0, 0.12), 0px 8px 16px -8px rgba(0, 0, 0, 0.16)"
  },
  "font": {
    "--el-font-size-extra-large": "20px",
    "--el-font-size-large": "18px",
    "--el-font-size-medium": "16px",
    "--el-font-size-base": "14px",
    "--el-font-size-small": "13px",
    "--el-font-size-extra-small": "12px",
    "--el-font-family": "'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif",
    "--el-font-weight-primary": "500",
    "--el-font-line-height-primary": "24px"
  },
  "size": {
    "--el-component-size-large": "40px",
    "--el-component-size": "32px",
    "--el-component-size-small": "24px"
  },
  "z-index": {
    "--el-index-normal": "1",
    "--el-index-top": "1000",
    "--el-index-popper": "2000"
  },
  "components": {
    "button": {
      "--el-button-font-weight": "var(--el-font-weight-primary)",
      "--el-button-border-color": "var(--el-border-color)",
      "--el-button-bg-color": "var(--el-fill-color-blank)",
      "--el-button-text-color": "var(--el-text-color-regular)",
      "--el-button-disabled-text-color": "var(--el-disabled-text-color)",
      "--el-button-disabled-bg-color": "var(--el-fill-color-blank)",
      "--el-button-disabled-border-color": "var(--el-border-color-light)",
      "--el-button-hover-text-color": "var(--el-color-primary)",
      "--el-button-hover-bg-color": "var(--el-color-primary-light-9)",
      "--el-button-hover-border-color": "var(--el-color-primary-light-7)",
      "--el-button-active-text-color": "var(--el-button-hover-text-color)",
      "--el-button-active-border-color": "var(--el-color-primary)",
      "--el-button-active-bg-color": "var(--el-button-hover-bg-color)",
      "--el-button-outline-color": "var(--el-color-primary-light-5)"
    },
    "input": {
      "--el-input-text-color": "var(--el-text-color-regular)",
      "--el-input-border": "var(--el-border)",
      "--el-input-hover-border": "var(--el-border-color-hover)",
      "--el-input-focus-border": "var(--el-color-primary)",
      "--el-input-transparent-border": "0 0 0 1px transparent inset",
      "--el-input-border-color": "var(--el-border-color)",
      "--el-input-border-radius": "var(--el-border-radius-base)",
      "--el-input-bg-color": "var(--el-fill-color-blank)",
      "--el-input-icon-color": "var(--el-text-color-placeholder)",
      "--el-input-placeholder-color": "var(--el-text-color-placeholder)",
      "--el-input-hover-border-color": "var(--el-border-color-hover)",
      "--el-input-clear-hover-color": "var(--el-text-color-secondary)",
      "--el-input-focus-border-color": "var(--el-color-primary)"
    },
    "tag": {
      "--el-tag-bg-color": "var(--el-color-primary-light-9)",
      "--el-tag-text-color": "var(--el-color-primary)",
      "--el-tag-border-color": "var(--el-color-primary-light-8)",
      "--el-tag-hover-color": "var(--el-color-primary)"
    },
    "notification": {
      "--el-notification-width": "330px",
      "--el-notification-padding": "14px 26px 14px 13px",
      "--el-notification-radius": "8px",
      "--el-notification-shadow": "var(--el-box-shadow-light)",
      "--el-notification-border-color": "var(--el-border-color-lighter)",
      "--el-notification-icon-size": "24px",
      "--el-notification-close-font-size": "16px",
      "--el-notification-group-margin-left": "13px",
      "--el-notification-group-margin-right": "8px",
      "--el-notification-content-font-size": "13px",
      "--el-notification-content-color": "var(--el-text-color-regular)",
      "--el-notification-title-font-size": "16px",
      "--el-notification-title-color": "var(--el-text-color-primary)",
      "--el-notification-close-color": "var(--el-text-color-secondary)",
      "--el-notification-close-hover-color": "var(--el-text-color-regular)"
    }
  }
}

const themes :{ [k: string]: any } = {
  'default': defaultTheme,
  'bubblegum': bubblegumTheme,
  'vintagePaper': vintagePaperTheme,
}

const applyTheme = function(data: any) {
  if (data instanceof Object) {
    Object.keys(data).forEach((key) => {
      if (data[key] instanceof Object) {
        applyTheme(data[key])
      } else {
        document.documentElement.style.setProperty(key, data[key])
      }
    })
  }
}

export default function setAsDarkTheme(darkMode: boolean) {
  const theme = window.localStorage.getItem('theme')
  if (theme) {
    applyTheme(themes[theme])
  }
  document.documentElement.className = darkMode ? 'dark' : 'light'
}
