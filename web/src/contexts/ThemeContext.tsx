import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import type { Theme, ThemeContextValue } from '../types';

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

const THEME_STORAGE_KEY = 'ui-theme';

interface ThemeProviderProps {
  children: ReactNode;
  defaultTheme?: Theme;
}

export function ThemeProvider({ children, defaultTheme = 'modern-blue' }: ThemeProviderProps) {
  const [theme, setThemeState] = useState<Theme>(() => {
    // Try to load theme from localStorage
    const stored = localStorage.getItem(THEME_STORAGE_KEY);
    return (stored as Theme) || defaultTheme;
  });

  useEffect(() => {
    // Remove all theme classes
    document.documentElement.classList.remove(
      'theme-modern-blue',
      'theme-warm-sunset',
      'theme-neo-mint',
      'theme-slate-dark',
      'theme-purple-dream',
      'theme-ocean-breeze',
      'theme-forest-green',
      'theme-rose-gold',
      'theme-midnight-purple',
      'theme-sakura-pink',
      'theme-cyber-neon'
    );

    // Add current theme class
    document.documentElement.classList.add(`theme-${theme}`);

    // Save to localStorage
    localStorage.setItem(THEME_STORAGE_KEY, theme);
  }, [theme]);

  const setTheme = (newTheme: Theme) => {
    setThemeState(newTheme);
  };

  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
