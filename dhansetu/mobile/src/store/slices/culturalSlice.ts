import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Quote {
  id: string;
  quote: string;
  author: string;
  language: string;
  category: string;
}

interface Festival {
  id: string;
  name: string;
  date: string;
  theme: {
    primaryColor: string;
    secondaryColor: string;
    accentColor: string;
  };
  greetings: string[];
}

interface CulturalState {
  currentQuote: Quote | null;
  currentFestival: Festival | null;
  quotes: Quote[];
  festivals: Festival[];
  selectedLanguage: string;
  patriotismScore: number;
}

const initialState: CulturalState = {
  currentQuote: null,
  currentFestival: null,
  quotes: [],
  festivals: [],
  selectedLanguage: 'en',
  patriotismScore: 0,
};

const culturalSlice = createSlice({
  name: 'cultural',
  initialState,
  reducers: {
    setCurrentQuote: (state, action: PayloadAction<Quote>) => {
      state.currentQuote = action.payload;
    },
    setCurrentFestival: (state, action: PayloadAction<Festival | null>) => {
      state.currentFestival = action.payload;
    },
    setQuotes: (state, action: PayloadAction<Quote[]>) => {
      state.quotes = action.payload;
    },
    setFestivals: (state, action: PayloadAction<Festival[]>) => {
      state.festivals = action.payload;
    },
    setSelectedLanguage: (state, action: PayloadAction<string>) => {
      state.selectedLanguage = action.payload;
    },
    incrementPatriotismScore: (state, action: PayloadAction<number>) => {
      state.patriotismScore += action.payload;
    },
    resetCultural: () => initialState,
  },
});

export const {
  setCurrentQuote,
  setCurrentFestival,
  setQuotes,
  setFestivals,
  setSelectedLanguage,
  incrementPatriotismScore,
  resetCultural,
} = culturalSlice.actions;

export default culturalSlice.reducer;