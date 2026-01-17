import React, { createContext, useContext, useEffect, useState } from 'react'

type Theme = 'light' | 'dark'

interface ThemeContextType {
  theme: Theme
  toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>(() => {
    try {
      // Check localStorage first
      const savedTheme = localStorage.getItem('theme')
      if (savedTheme === 'light' || savedTheme === 'dark') {
        return savedTheme
      }
    } catch (e) {
      console.warn('[ThemeContext] Failed to access localStorage:', e)
    }
    
    // Fallback to system preference
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark'
    }
    return 'light'
  })

  useEffect(() => {
    const root = window.document.documentElement
    console.log('[ThemeContext] Theme effect triggered. Current theme:', theme)
    
    const oldClass = root.classList.contains('dark') ? 'dark' : 'light'
    root.classList.remove('light', 'dark')
    root.classList.add(theme)
    
    console.log(`[ThemeContext] Updated HTML class from '${oldClass}' to '${theme}'`)
    try {
      localStorage.setItem('theme', theme)
    } catch (e) {
      console.warn('[ThemeContext] Failed to write localStorage:', e)
    }
  }, [theme])

  const toggleTheme = () => {
    console.log('[ThemeContext] toggleTheme called. Previous theme:', theme)
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'))
  }

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}

export function useTheme() {
  const context = useContext(ThemeContext)
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider')
  }
  return context
}
