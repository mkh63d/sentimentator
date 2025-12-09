import { createI18n } from 'vue-i18n';
import en from './locales/en.json';
import pl from './locales/pl.json';

const messages = {
  en,
  pl,
};

const savedLocale = localStorage.getItem('locale') || 'en';

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'en',
  messages,
  globalInjection: true,
});

export default i18n;
