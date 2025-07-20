"use strict";

import NativeModule from "./NativeRNLocalize";
export function getCalendar() {
  return NativeModule.getCalendar();
}
export function getCountry() {
  return NativeModule.getCountry();
}
export function getCurrencies() {
  return NativeModule.getCurrencies();
}
export function getLocales() {
  return NativeModule.getLocales();
}
export function getNumberFormatSettings() {
  return NativeModule.getNumberFormatSettings();
}
export function getTemperatureUnit() {
  return NativeModule.getTemperatureUnit();
}
export function getTimeZone() {
  return NativeModule.getTimeZone();
}
export function uses24HourClock() {
  return NativeModule.uses24HourClock();
}
export function usesMetricSystem() {
  return NativeModule.usesMetricSystem();
}
export function usesAutoDateAndTime() {
  return NativeModule.usesAutoDateAndTime() ?? undefined;
}
export function usesAutoTimeZone() {
  return NativeModule.usesAutoTimeZone() ?? undefined;
}
export async function openAppLanguageSettings() {
  await NativeModule.openAppLanguageSettings();
}
//# sourceMappingURL=module.js.map