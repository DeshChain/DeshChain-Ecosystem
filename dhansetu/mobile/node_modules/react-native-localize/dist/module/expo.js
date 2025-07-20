"use strict";

import * as Expo from "@expo/config-plugins";
import { mergeContents } from "@expo/config-plugins/build/utils/generateCode";
import { mkdirSync, writeFileSync } from "fs";
import { join } from "path";
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
    const xmlDir = join(config.modRequest.platformProjectRoot, "app", "src", "main", "res", "xml");
    mkdirSync(xmlDir, {
      recursive: true
    });
    writeFileSync(join(xmlDir, "locales_config.xml"), `<?xml version="1.0" encoding="utf-8"?>
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
    } = mergeContents({
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
export default Expo.createRunOncePlugin(withAppLocales, PACKAGE_NAME);
//# sourceMappingURL=expo.js.map