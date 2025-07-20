# DeshChain Money Order - 22 Language Support Implementation

This directory contains the internationalization (i18n) implementation for supporting all 22 official languages of India in the DeshChain Money Order system.

## Supported Languages

| Language | Code | Native Name | Script | Direction | Status |
|----------|------|-------------|---------|-----------|---------|
| English | en | English | Latin | LTR | ✅ Complete |
| Hindi | hi | हिन्दी | Devanagari | LTR | ✅ Complete |
| Bengali | bn | বাংলা | Bengali | LTR | ✅ Complete |
| Telugu | te | తెలుగు | Telugu | LTR | 🟡 Basic |
| Marathi | mr | मराठी | Devanagari | LTR | 🟡 Basic |
| Tamil | ta | தமிழ் | Tamil | LTR | 🟡 Basic |
| Urdu | ur | اردو | Arabic | RTL | 🟡 Basic |
| Gujarati | gu | ગુજરાતી | Gujarati | LTR | 🟡 Basic |
| Kannada | kn | ಕನ್ನಡ | Kannada | LTR | 🟡 Basic |
| Odia | or | ଓଡ଼ିଆ | Odia | LTR | 🔴 Placeholder |
| Malayalam | ml | മലയാളം | Malayalam | LTR | 🔴 Placeholder |
| Punjabi | pa | ਪੰਜਾਬੀ | Gurmukhi | LTR | 🔴 Placeholder |
| Assamese | as | অসমীয়া | Bengali | LTR | 🔴 Placeholder |
| Maithili | mai | मैथिली | Devanagari | LTR | 🔴 Placeholder |
| Santali | sat | ᱥᱟᱱᱛᱟᱲᱤ | Ol Chiki | LTR | 🔴 Placeholder |
| Kashmiri | ks | کٲشُر | Arabic | RTL | 🔴 Placeholder |
| Nepali | ne | नेपाली | Devanagari | LTR | 🔴 Placeholder |
| Konkani | kok | कोंकणी | Devanagari | LTR | 🔴 Placeholder |
| Sindhi | sd | सिन्धी | Devanagari | LTR | 🔴 Placeholder |
| Dogri | doi | डोगरी | Devanagari | LTR | 🔴 Placeholder |
| Manipuri | mni | মৈতৈলোন্ | Bengali | LTR | 🔴 Placeholder |
| Sanskrit | sa | संस्कृतम् | Devanagari | LTR | 🔴 Placeholder |

## Features

### 1. Complete Language Support
- All 22 official Indian languages
- RTL support for Urdu and Kashmiri
- Script-specific font loading
- Localized number systems

### 2. Cultural Context
- Time-based greetings
- Cultural phrases and blessings
- Festival-specific themes
- Regional formatting preferences

### 3. Number Localization
- Script-specific numerals (०१२३४५६७८९ for Devanagari)
- Indian numbering system (lakhs and crores)
- Currency formatting
- Phone number validation

### 4. User Experience
- Language selector with search
- Persistent language preference
- Automatic font loading
- RTL layout support

## Usage

### Basic Usage

```typescript
import { useTranslation } from 'react-i18next';
import { useLanguage } from '../hooks/useLanguage';

function MyComponent() {
  const { t } = useTranslation();
  const { formatCurrency, currentLanguage } = useLanguage();
  
  return (
    <div>
      <h1>{t('app.title')}</h1>
      <p>{formatCurrency(1000)}</p> {/* ₹1,000 */}
    </div>
  );
}
```

### Language Selector

```typescript
import { LanguageSelector } from '../components/LanguageSelector';

function Header() {
  return (
    <header>
      <LanguageSelector 
        variant="button"
        showNativeName={true}
        showRegion={true}
      />
    </header>
  );
}
```

### RTL Support

```typescript
import { useRTL } from '../hooks/useLanguage';

function RTLComponent() {
  const { isRTL, direction, textAlign } = useRTL();
  
  return (
    <div dir={direction} style={{ textAlign }}>
      {/* Content automatically adjusts for RTL languages */}
    </div>
  );
}
```

### Number Localization

```typescript
import { useLocalizedNumbers } from '../hooks/useLanguage';

function NumberInput() {
  const { localizeInput, formatForDisplay } = useLocalizedNumbers();
  const [value, setValue] = useState('');
  
  const handleChange = (e) => {
    // Convert localized digits to Western for processing
    const normalized = localizeInput(e.target.value);
    setValue(normalized);
  };
  
  return (
    <input
      value={formatForDisplay(value)} // Display with localized digits
      onChange={handleChange}
    />
  );
}
```

## File Structure

```
i18n/
├── index.ts                 # Main i18n configuration
├── LanguageService.ts       # Language utilities and helpers
├── locales/                 # Translation files
│   ├── en.json             # English (complete)
│   ├── hi.json             # Hindi (complete)
│   ├── bn.json             # Bengali (complete)
│   └── ...                 # Other languages
├── translationGenerator.ts  # Tool for generating translations
└── README.md               # This file
```

## Adding Translations

### 1. Update Translation File

Edit the appropriate JSON file in `locales/`:

```json
{
  "moneyOrder": {
    "form": {
      "amount": "राशि"  // Hindi translation
    }
  }
}
```

### 2. Add Cultural Context

Update `LanguageService.ts` with cultural phrases:

```typescript
const phrases = {
  hi: {
    thanks: 'धन्यवाद 🙏',
    welcome: 'आपका स्वागत है',
    blessing: 'जय हिंद! आपकी यात्रा शुभ हो'
  }
};
```

### 3. Test RTL Languages

For Urdu and Kashmiri, ensure proper RTL handling:

```typescript
// Automatically handled by the system
document.documentElement.dir = 'rtl';
```

## Font Requirements

The system automatically loads appropriate fonts:

- **Devanagari**: Noto Sans Devanagari (Hindi, Marathi, Nepali, etc.)
- **Bengali**: Noto Sans Bengali (Bengali, Assamese)
- **Tamil**: Noto Sans Tamil
- **Arabic**: Noto Nastaliq Urdu (Urdu, Kashmiri)
- **Ol Chiki**: Noto Sans Ol Chiki (Santali)

## Best Practices

1. **Always use translation keys**: Never hardcode text
2. **Provide context**: Use descriptive translation keys
3. **Test with longest translations**: German/Sanskrit often longest
4. **Validate phone numbers**: Use region-specific validation
5. **Format numbers properly**: Use Indian numbering system
6. **Consider cultural context**: Add appropriate greetings/phrases

## Contribution Guidelines

1. **Complete translations**: Focus on completing one language at a time
2. **Native speakers**: Translations should be done by native speakers
3. **Cultural sensitivity**: Ensure appropriate cultural context
4. **Consistency**: Maintain consistent terminology across languages
5. **Testing**: Test with actual users from the language community

## Future Enhancements

- [ ] Voice input in local languages
- [ ] Transliteration support
- [ ] Regional calendar systems
- [ ] Local payment terminology
- [ ] Speech-to-text in all languages
- [ ] AI-powered translation suggestions

## Resources

- [Official Languages of India](https://en.wikipedia.org/wiki/Languages_with_official_status_in_India)
- [Indian Numbering System](https://en.wikipedia.org/wiki/Indian_numbering_system)
- [Unicode Scripts](https://www.unicode.org/charts/)
- [Google Noto Fonts](https://www.google.com/get/noto/)