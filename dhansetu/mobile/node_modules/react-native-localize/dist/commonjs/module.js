"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.getCalendar = getCalendar;
exports.getCountry = getCountry;
exports.getCurrencies = getCurrencies;
exports.getLocales = getLocales;
exports.getNumberFormatSettings = getNumberFormatSettings;
exports.getTemperatureUnit = getTemperatureUnit;
exports.getTimeZone = getTimeZone;
exports.openAppLanguageSettings = openAppLanguageSettings;
exports.uses24HourClock = uses24HourClock;
exports.usesAutoDateAndTime = usesAutoDateAndTime;
exports.usesAutoTimeZone = usesAutoTimeZone;
exports.usesMetricSystem = usesMetricSystem;
var _NativeRNLocalize = _interopRequireDefault(require("./NativeRNLocalize"));
function _interopRequireDefault(e) { return e && e.__esModule ? e : { default: e }; }
function getCalendar() {
  return _NativeRNLocalize.default.getCalendar();
}
function getCountry() {
  return _NativeRNLocalize.default.getCountry();
}
function getCurrencies() {
  return _NativeRNLocalize.default.getCurrencies();
}
function getLocales() {
  return _NativeRNLocalize.default.getLocales();
}
function getNumberFormatSettings() {
  return _NativeRNLocalize.default.getNumberFormatSettings();
}
function getTemperatureUnit() {
  return _NativeRNLocalize.default.getTemperatureUnit();
}
function getTimeZone() {
  return _NativeRNLocalize.default.getTimeZone();
}
function uses24HourClock() {
  return _NativeRNLocalize.default.uses24HourClock();
}
function usesMetricSystem() {
  return _NativeRNLocalize.default.usesMetricSystem();
}
function usesAutoDateAndTime() {
  return _NativeRNLocalize.default.usesAutoDateAndTime() ?? undefined;
}
function usesAutoTimeZone() {
  return _NativeRNLocalize.default.usesAutoTimeZone() ?? undefined;
}
async function openAppLanguageSettings() {
  await _NativeRNLocalize.default.openAppLanguageSettings();
}
//# sourceMappingURL=module.js.map