import React, { createContext, useContext, useState, ReactNode } from 'react'
import { Bus } from '../types/bus'
import axios from 'axios'
import { useAuth } from './AuthContext'

interface HttpContextProps {
  busList: Bus[]
  loadBusList: (s: string) => Promise<void>
}

const HttpContext = createContext<HttpContextProps | undefined>(undefined)

// eslint-disable-next-line react-refresh/only-export-components
export const useHttp = () => {
  const context = useContext(HttpContext)

  if (!context) {
    throw new Error('useHttp must be used within a Http Provider')
  }

  return context
}

interface HttpProviderProps {
  children: ReactNode
}

const httpUrl = import.meta.env.VITE_APP_HTTP_SERVER_LINK

export const HttpProvider: React.FC<HttpProviderProps> = ({ children }) => {
  const [busList, setBusList] = useState<Bus[]>([])
  const { session } = useAuth()

  const loadBusList = async (s: string) => {
    try {
      const res = await axios.get(`${httpUrl}/bus?search=${s}`, {
        headers: {
          Authorization: 'bearer ' + session?.access_token
        }
      })
      setBusList(res.data.buses)
    } catch (e) {
      console.log(e)
    }
  }

  const contextValue: HttpContextProps = {
    busList,
    loadBusList
  }
  return (
    <HttpContext.Provider value={contextValue}>{children}</HttpContext.Provider>
  )
}
