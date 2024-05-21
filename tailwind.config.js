/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./src/components/*.go",
        "./src/components/*.templ",
    ],
    theme: {
        extend: {},
    },
    plugins: [],
    safelist: ["text-gray-500"]
}

