import * as React from 'react'
import { createContext, useContext, useEffect, useState } from 'react'

type Theme = 'light' | 'dark'

interface ThemeContextType {
  theme: Theme
  toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

/**
 * Get initial theme safely
 */
function getInitialTheme(): Theme {
  try {
    const savedTheme = typeof window !== 'undefined' ? localStorage.getItem('theme') : null
    if (savedTheme === 'light' || savedTheme === 'dark') {
      return savedTheme
    }
    
    if (typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark'
    }
  } catch (e) {
    console.warn('[ThemeContext] Failed to determine initial theme:', e)
  }
  return 'light'
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>(getInitialTheme)

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
