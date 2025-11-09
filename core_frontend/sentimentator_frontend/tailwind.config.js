/** @type {import('tailwindcss').Config} */
export const content = ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'];
export const theme = {
  extend: {
    colors: {
      primary: '(var(--color-primary))',
      secondary: '(var(--color-secondary))',
      accent: '(var(--color-accent))',
      background: '(var(--color-background))',
      textBase: '(var(--color-text-base))',
      fourthary: '(var(--color-fourthary))',
    },
  },
  darkMode: ['class', '[data-theme="dark"]'],
};
export const plugins = [];
