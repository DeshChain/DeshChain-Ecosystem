import React, { createContext, useContext, useState, useEffect } from 'react';
import { I18n } from 'i18n-js';
import * as Localization from 'react-native-localize';

interface LocalizationContextType {
  i18n: I18n;
  currentLocale: string;
  changeLanguage: (locale: string) => void;
}

const LocalizationContext = createContext<LocalizationContextType | undefined>(undefined);

export const useLocalization = () => {
  const context = useContext(LocalizationContext);
  if (!context) {
    throw new Error('useLocalization must be used within a LocalizationProvider');
  }
  return context;
};

export const LocalizationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [i18n] = useState(new I18n({
    en: {
      welcome: "Welcome to DhanSetu",
      wallet: "Wallet",
      send: "Send",
      receive: "Receive"
    },
    hi: {
      welcome: "धनसेतु में आपका स्वागत है",
      wallet: "बटुआ",
      send: "भेजें",
      receive: "प्राप्त करें"
    }
  }));
  
  const [currentLocale, setCurrentLocale] = useState('en');

  useEffect(() => {
    const locales = Localization.getLocales();
    if (locales[0]) {
      const locale = locales[0].languageCode;
      setCurrentLocale(locale);
      i18n.locale = locale;
    }
  }, [i18n]);

  const changeLanguage = (locale: string) => {
    setCurrentLocale(locale);
    i18n.locale = locale;
  };

  return (
    <LocalizationContext.Provider value={{ i18n, currentLocale, changeLanguage }}>
      {children}
    </LocalizationContext.Provider>
  );
};