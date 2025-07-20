/*
Copyright 2024 DeshChain Foundation

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

import fs from 'fs';
import path from 'path';

// Basic translations for all 22 languages
const basicTranslations: Record<string, any> = {
  te: {
    app: {
      title: "దేశ్‌చైన్ మనీ ఆర్డర్",
      tagline: "డిజిటల్ భారత్ కోసం డిజిటల్ మనీ ఆర్డర్లు",
      description: "బ్లాక్‌చైన్ భద్రతతో భారతదేశం అంతటా తక్షణమే డబ్బును పంపండి"
    },
    navigation: {
      home: "హోమ్",
      send: "డబ్బు పంపండి",
      receive: "స్వీకరించండి",
      track: "ఆర్డర్ ట్రాక్ చేయండి",
      history: "చరిత్ర",
      help: "సహాయం",
      settings: "సెట్టింగులు",
      logout: "లాగ్అవుట్"
    },
    common: {
      loading: "లోడ్ అవుతోంది...",
      error: "లోపం",
      success: "విజయవంతం",
      confirm: "నిర్ధారించండి",
      cancel: "రద్దు చేయండి",
      yes: "అవును",
      no: "కాదు"
    }
  },
  mr: {
    app: {
      title: "देशचेन मनी ऑर्डर",
      tagline: "डिजिटल भारतासाठी डिजिटल मनी ऑर्डर",
      description: "ब्लॉकचेन सुरक्षेसह संपूर्ण भारतात त्वरित पैसे पाठवा"
    },
    navigation: {
      home: "मुख्यपृष्ठ",
      send: "पैसे पाठवा",
      receive: "प्राप्त करा",
      track: "ऑर्डर ट्रॅक करा",
      history: "इतिहास",
      help: "मदत",
      settings: "सेटिंग्ज",
      logout: "लॉगआउट"
    },
    common: {
      loading: "लोड होत आहे...",
      error: "त्रुटी",
      success: "यशस्वी",
      confirm: "पुष्टी करा",
      cancel: "रद्द करा",
      yes: "होय",
      no: "नाही"
    }
  },
  ta: {
    app: {
      title: "தேஷ்செயின் பணக்கட்டளை",
      tagline: "டிஜிட்டல் இந்தியாவிற்கான டிஜிட்டல் பணக்கட்டளைகள்",
      description: "பிளாக்செயின் பாதுகாப்புடன் இந்தியா முழுவதும் உடனடியாக பணம் அனுப்புங்கள்"
    },
    navigation: {
      home: "முகப்பு",
      send: "பணம் அனுப்பு",
      receive: "பெறு",
      track: "ஆர்டரை கண்காணி",
      history: "வரலாறு",
      help: "உதவி",
      settings: "அமைப்புகள்",
      logout: "வெளியேறு"
    },
    common: {
      loading: "ஏற்றுகிறது...",
      error: "பிழை",
      success: "வெற்றி",
      confirm: "உறுதிப்படுத்து",
      cancel: "ரத்து செய்",
      yes: "ஆம்",
      no: "இல்லை"
    }
  },
  ur: {
    app: {
      title: "دیش چین منی آرڈر",
      tagline: "ڈیجیٹل بھارت کے لیے ڈیجیٹل منی آرڈر",
      description: "بلاک چین سیکیورٹی کے ساتھ پورے بھارت میں فوری طور پر رقم بھیجیں"
    },
    navigation: {
      home: "ہوم",
      send: "رقم بھیجیں",
      receive: "وصول کریں",
      track: "آرڈر ٹریک کریں",
      history: "تاریخ",
      help: "مدد",
      settings: "ترتیبات",
      logout: "لاگ آؤٹ"
    },
    common: {
      loading: "لوڈ ہو رہا ہے...",
      error: "خرابی",
      success: "کامیاب",
      confirm: "تصدیق کریں",
      cancel: "منسوخ کریں",
      yes: "جی ہاں",
      no: "نہیں"
    }
  },
  gu: {
    app: {
      title: "દેશચેન મની ઓર્ડર",
      tagline: "ડિજિટલ ભારત માટે ડિજિટલ મની ઓર્ડર",
      description: "બ્લોકચેન સુરક્ષા સાથે સમગ્ર ભારતમાં તરત જ પૈસા મોકલો"
    },
    navigation: {
      home: "હોમ",
      send: "પૈસા મોકલો",
      receive: "પ્રાપ્ત કરો",
      track: "ઓર્ડર ટ્રેક કરો",
      history: "ઇતિહાસ",
      help: "મદદ",
      settings: "સેટિંગ્સ",
      logout: "લોગઆઉટ"
    },
    common: {
      loading: "લોડ થઈ રહ્યું છે...",
      error: "ભૂલ",
      success: "સફળ",
      confirm: "પુષ્ટિ કરો",
      cancel: "રદ કરો",
      yes: "હા",
      no: "ના"
    }
  },
  kn: {
    app: {
      title: "ದೇಶ್‌ಚೈನ್ ಮನಿ ಆರ್ಡರ್",
      tagline: "ಡಿಜಿಟಲ್ ಭಾರತಕ್ಕೆ ಡಿಜಿಟಲ್ ಮನಿ ಆರ್ಡರ್‌ಗಳು",
      description: "ಬ್ಲಾಕ್‌ಚೈನ್ ಭದ್ರತೆಯೊಂದಿಗೆ ಭಾರತದಾದ್ಯಂತ ತಕ್ಷಣ ಹಣವನ್ನು ಕಳುಹಿಸಿ"
    },
    navigation: {
      home: "ಮುಖಪುಟ",
      send: "ಹಣ ಕಳುಹಿಸಿ",
      receive: "ಸ್ವೀಕರಿಸಿ",
      track: "ಆರ್ಡರ್ ಟ್ರ್ಯಾಕ್ ಮಾಡಿ",
      history: "ಇತಿಹಾಸ",
      help: "ಸಹಾಯ",
      settings: "ಸೆಟ್ಟಿಂಗ್‌ಗಳು",
      logout: "ಲಾಗ್ಔಟ್"
    },
    common: {
      loading: "ಲೋಡ್ ಆಗುತ್ತಿದೆ...",
      error: "ದೋಷ",
      success: "ಯಶಸ್ವಿ",
      confirm: "ದೃಢೀಕರಿಸಿ",
      cancel: "ರದ್ದುಮಾಡಿ",
      yes: "ಹೌದು",
      no: "ಇಲ್ಲ"
    }
  }
};

// Generate placeholder structure for remaining languages
const generatePlaceholderStructure = (languageCode: string, nativeName: string): any => {
  return {
    app: {
      title: `DeshChain Money Order (${nativeName})`,
      tagline: `Digital Money Orders for Digital India (${nativeName})`,
      description: `Send money across India instantly with blockchain security (${nativeName})`
    },
    navigation: {
      home: "Home",
      send: "Send Money",
      receive: "Receive",
      track: "Track Order",
      history: "History",
      help: "Help",
      settings: "Settings",
      logout: "Logout"
    },
    common: {
      loading: "Loading...",
      error: "Error",
      success: "Success",
      confirm: "Confirm",
      cancel: "Cancel",
      yes: "Yes",
      no: "No"
    }
  };
};

// Languages that need placeholder files
const placeholderLanguages = [
  { code: 'or', name: 'Odia', nativeName: 'ଓଡ଼ିଆ' },
  { code: 'ml', name: 'Malayalam', nativeName: 'മലയാളം' },
  { code: 'pa', name: 'Punjabi', nativeName: 'ਪੰਜਾਬੀ' },
  { code: 'as', name: 'Assamese', nativeName: 'অসমীয়া' },
  { code: 'mai', name: 'Maithili', nativeName: 'मैथिली' },
  { code: 'sat', name: 'Santali', nativeName: 'ᱥᱟᱱᱛᱟᱲᱤ' },
  { code: 'ks', name: 'Kashmiri', nativeName: 'کٲشُر' },
  { code: 'ne', name: 'Nepali', nativeName: 'नेपाली' },
  { code: 'kok', name: 'Konkani', nativeName: 'कोंकणी' },
  { code: 'sd', name: 'Sindhi', nativeName: 'सिन्धी' },
  { code: 'doi', name: 'Dogri', nativeName: 'डोगरी' },
  { code: 'mni', name: 'Manipuri', nativeName: 'মৈতৈলোন্' },
  { code: 'sa', name: 'Sanskrit', nativeName: 'संस्कृतम्' }
];

// Export function to generate all translation files
export const generateTranslationFiles = (outputDir: string) => {
  // Create basic translations
  Object.entries(basicTranslations).forEach(([langCode, translations]) => {
    const filePath = path.join(outputDir, `${langCode}.json`);
    fs.writeFileSync(filePath, JSON.stringify(translations, null, 2));
    console.log(`Created translation file: ${filePath}`);
  });

  // Create placeholder translations
  placeholderLanguages.forEach(({ code, name, nativeName }) => {
    const translations = generatePlaceholderStructure(code, nativeName);
    const filePath = path.join(outputDir, `${code}.json`);
    fs.writeFileSync(filePath, JSON.stringify(translations, null, 2));
    console.log(`Created placeholder translation file: ${filePath}`);
  });
};

// CLI usage
if (require.main === module) {
  const outputDir = path.join(__dirname, 'locales');
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }
  generateTranslationFiles(outputDir);
}