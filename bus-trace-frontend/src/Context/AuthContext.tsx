import React, {
  createContext,
  useContext,
  useState,
  ReactNode,
  useEffect
} from 'react'
import { supabase } from '../helpers/Supabase'
import { Session } from '@supabase/supabase-js'
import { User } from '../types/user'
import axios from 'axios'

interface AuthContextProps {
  authenticated: boolean
  session: Session | null
  userData: User | null
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
  const [userData, setUserData] = useState<User | null>(null)

  const contextValue: AuthContextProps = {
    authenticated,
    session,
    userData
  }

  useEffect(() => {
    supabase.auth.onAuthStateChange(async (_, session) => {
      if (session) {
        const res = await axios.get(
          `${import.meta.env.VITE_APP_HTTP_SERVER_LINK}/user`,
          {
            headers: {
              Authorization: 'bearer ' + session?.access_token
            }
          }
        )
        setUserData(res.data.user)
        setAuthenticated(true)
        setSession(session)
      } else {
        setUserData(null)
        setAuthenticated(false)
        setSession(null)
      }
    })
  }, [])

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  )
}
