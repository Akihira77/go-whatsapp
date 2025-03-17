/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.templ",
    "./views/**/*.go",
    "./**/*.templ",
    "./**/*.go",
  ],
  theme: {
    extend: {
      extend: {
        colors: {
          whatsapp: {
            light: "#25D366",
            DEFAULT: "#128C7E",
            dark: "#075E54",
          },
        },
      },
    },
  },
  plugins: [],
};
