import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import HttpApi from 'i18next-http-backend';

export const SupportedLanguages = ['ar', 'de', 'en', 'es', 'fr', 'ja', 'ru'];

i18n
	.use(HttpApi)
	.use(initReactI18next)
	.init({
		debug: false,
		lng: 'en',
		supportedLngs: SupportedLanguages,
		fallbackLng: 'en',
		interpolation: {
			escapeValue: false, // not needed for react as it escapes by default
		},
		backend: {
			loadPath: process.env.PUBLIC_URL+'/locales/{{lng}}/{{ns}}.json',
		}
	});

export default i18n;
