import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

// Debug: Monitor HTML class changes
const observer = new MutationObserver((mutations) => {
  mutations.forEach((mutation) => {
    if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
      console.log('[DOM Monitor] HTML class changed to:', document.documentElement.className)
    }
  })
})

observer.observe(document.documentElement, {
  attributes: true,
  attributeFilter: ['class']
})

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
