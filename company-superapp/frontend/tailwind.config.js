/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./App.{js,jsx,ts,tsx}", "./src/**/*.{js,jsx,ts,tsx}"],
    theme: {
        extend: {
            colors: {
                primary: '#FF4B33',
                background: '#0F0F0F',
                surface: '#1E1E1E',
            }
        },
    },
    plugins: [],
}
