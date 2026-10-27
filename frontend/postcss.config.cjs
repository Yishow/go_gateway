module.exports = {
  plugins: {
    '@tailwindcss/postcss': {
      // Reparse generated keyframes with source metadata before Vite rewrites URLs.
      optimize: { minify: process.env.NODE_ENV === 'production' },
    },
  },
}
