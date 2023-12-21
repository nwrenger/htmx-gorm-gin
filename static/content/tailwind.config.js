/** @type {import('tailwindcss').Config} */
export const content = ["./../../web/**/*.{go,tmpl}"];
export const theme = {
  extend: {},
};
export const plugins = [
  require('@tailwindcss/forms'),
  require('daisyui'),
];