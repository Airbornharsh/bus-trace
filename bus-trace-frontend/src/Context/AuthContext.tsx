import React, {
  createContext,
  useContext,
  useState,
  ReactNode,
  useEffect
} from 'react'
import { supabase } from '../helpers/Supabase'
import { User } from '@supabase/supabase-js'

interface AuthContextProps {
  authenticated: boolean
  user: User | null
}

const AuthContext = createContext<AuthContextProps | undefined>(undefined)

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used within a AuthProvider')
  }

  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState<User | null>(null)

  const contextValue: AuthContextProps = {
    authenticated,
    user
  }

  useEffect(() => {
    supabase.auth.onAuthStateChange((event, session) => {
      console.log('Event:', event)
      console.log('Session:', session)
      if (event === 'SIGNED_IN') {
        setAuthenticated(true)
        setUser(session?.user || null)
      } else {
        setAuthenticated(false)
        setUser(null)
      }
    })
  }, [])

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  )
}
