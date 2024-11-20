import translations from '$lib/translations.json';

export function translate(key: string, language: string = "pt") : string {
    return translations[key] || key;
}