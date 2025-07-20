"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;
var Expo = _interopRequireWildcard(require("@expo/config-plugins"));
var _generateCode = require("@expo/config-plugins/build/utils/generateCode");
var _fs = require("fs");
var _path = require("path");
function _interopRequireWildcard(e, t) { if ("function" == typeof WeakMap) var r = new WeakMap(), n = new WeakMap(); return (_interopRequireWildcard = function (e, t) { if (!t && e && e.__esModule) return e; var o, i, f = { __proto__: null, default: e }; if (null === e || "object" != typeof e && "function" != typeof e) return f; if (o = t ? n : r) { if (o.has(e)) return o.get(e); o.set(e, f); } for (const t in e) "default" !== t && {}.hasOwnProperty.call(e, t) && ((i = (o = Object.defineProperty) && Object.getOwnPropertyDescriptor(e, t)) && (i.get || i.set) ? o(f, t, i) : f[t] = e[t]); return f; })(e, t); }
const PACKAGE_NAME = "react-native-localize";
const withIosAppLocales = ({
  ios = {},
  ...config
}, locales) => ({
  ...config,
  ios: {
    ...ios,
    infoPlist: {
      CFBundleLocalizations: locales
    }
  }
});
const withAndroidAppLocales = (config, locales) => {
  const withLocalesConfig = Expo.withDangerousMod(config, ["android", config => {
    const xmlDir = (0, _path.join)(config.modRequest.platformProjectRoot, "app", "src", "main", "res", "xml");
    (0, _fs.mkdirSync)(xmlDir, {
      recursive: true
    });
    (0, _fs.writeFileSync)((0, _path.join)(xmlDir, "locales_config.xml"), `<?xml version="1.0" encoding="utf-8"?>
<locale-config xmlns:android="http://schemas.android.com/apk/res/android">
${locales.map(locale => `  <locale android:name="${locale}"/>`).join("\n")}
</locale-config>
`);
    return config;
  }]);
  const withMainApplication = Expo.withAndroidManifest(withLocalesConfig, config => {
    const mainApplication = Expo.AndroidConfig.Manifest.getMainApplicationOrThrow(config.modResults);
    mainApplication.$ = {
      ...mainApplication.$,
      "android:localeConfig": "@xml/locales_config"
    };
    return config;
  });
  return Expo.withAppBuildGradle(withMainApplication, config => {
    const {
      modResults
    } = config;
    const list = locales.map(lang => `"${lang}"`).join(", ");
    const {
      contents
    } = (0, _generateCode.mergeContents)({
      src: modResults.contents,
      comment: "//",
      tag: PACKAGE_NAME,
      offset: 1,
      anchor: /versionName "[\d.]+"/,
      newSrc: `        resourceConfigurations += ${modResults.language === "kt" ? `listOf(${list})` : `[${list}]`}`
    });
    return {
      ...config,
      modResults: {
        ...modResults,
        contents
      }
    };
  });
};
const withAppLocales = (config, {
  locales = []
} = {}) => {
  const plugins = [];
  const {
    platforms = []
  } = config;
  const {
    android,
    ios
  } = Array.isArray(locales) ? {
    android: locales,
    ios: locales
  } : {
    android: locales.android ?? [],
    ios: locales.ios ?? []
  };
  if (platforms.includes("ios") && ios.length > 0) {
    plugins.push([withIosAppLocales, ios]);
  }
  if (platforms.includes("android") && android.length > 0) {
    plugins.push([withAndroidAppLocales, android]);
  }
  return Expo.withPlugins(config, plugins);
};
var _default = exports.default = Expo.createRunOncePlugin(withAppLocales, PACKAGE_NAME);
//# sourceMappingURL=expo.js.map