import React, {
  createContext,
  useContext,
  useState,
  ReactNode,
  useEffect
} from 'react'
import { supabase } from '../helpers/Supabase'
import { Session } from '@supabase/supabase-js'

interface AuthContextProps {
  authenticated: boolean
  session: Session | null
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
  const [session, setSession] = useState<Session | null>(null)

  const contextValue: AuthContextProps = {
    authenticated,
    session
  }

  useEffect(() => {
    supabase.auth.onAuthStateChange((event, session) => {
      if (
        event === 'SIGNED_IN' ||
        event === 'USER_UPDATED' ||
        event === 'INITIAL_SESSION' ||
        event === 'PASSWORD_RECOVERY' ||
        event === 'TOKEN_REFRESHED'
      ) {
        setAuthenticated(true)
        setSession(session)
      } else {
        setAuthenticated(false)
        setSession(null)
      }
    })
  }, [])

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  )
}
