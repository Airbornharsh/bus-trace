import React, { createContext, useContext, useState, ReactNode } from 'react'
import { Bus } from '../types/bus'
import axios from 'axios'
import { useAuth } from './AuthContext'

interface HttpContextProps {
  busList: Bus[]
  loadBusList: (s: string, first: boolean) => Promise<void>
  signUp: (
    id: string,
    name: string,
    email: string,
    phone: string
  ) => Promise<void>
  userDatas: {
    [key: string]: {
      name: string
      email: string
      phone: string
      lat: number
      long: number
      location: string
    }
  }
  userList: string[]
  setUserList: (users: string[]) => void
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
  const [userDatas, setUserDatas] = useState<{
    [key: string]: {
      name: string
      email: string
      phone: string
      lat: number
      long: number
      location: string
    }
  }>({})
  const [userList, setUserList] = useState<string[]>([])
  const [searchTimer, setSearchTimer] = useState<NodeJS.Timeout | null>(null)
  const { session } = useAuth()

  const loadBusList = async (s: string, first: boolean) => {
    if (searchTimer) {
      clearTimeout(searchTimer)
    }
    await new Promise((resolve, reject) => {
      const timer = setTimeout(
        () => {
          if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
              async (position) => {
                try {
                  if (
                    position.coords.latitude === 0 &&
                    position.coords.longitude === 0
                  ) {
                    throw new Error('Location not found')
                  }
                  const res = await axios.get(
                    `${httpUrl}/bus/${position.coords.latitude + 0.001}/${position.coords.longitude}?search=${s}&distance=5000`,
                    {
                      headers: {
                        Authorization: 'bearer ' + session?.access_token
                      }
                    }
                  )
                  setBusList(res.data.buses)
                  resolve(true)
                  return
                } catch (e) {
                  console.log(e)
                  reject(e)
                  return
                }
              },
              (e) => {
                console.error(e)
                reject(e)
                return
              },
              { enableHighAccuracy: true }
            )
          } else {
            console.log('Not Supported by Browser')
            alert('Not supported')
            reject('Not supported')
            return
          }
        },
        first ? 0 : 800
      )
      setSearchTimer(timer)
    })
  }

  const signUp = async (
    id: string,
    name: string,
    email: string,
    phone: string
  ) => {
    try {
      await axios.post(
        `${httpUrl}/user`,
        {
          id,
          name,
          email,
          phone
        },
        {
          headers: {
            Authorization: 'bearer ' + session?.access_token
          }
        }
      )
      return
    } catch (e) {
      console.log(e)
    }
  }

  const setUserListFn = async (users: string[]) => {
    users.forEach(async (u) => {
      try {
        if (Object.keys(userDatas).indexOf(u) === -1) {
          const res = await axios.get(`${httpUrl}/user/${u}`, {
            headers: {
              Authorization: 'bearer ' + session?.access_token
            }
          })
          console.log('User:', res.data.user)
          setUserDatas((prev) => ({
            ...prev,
            [u]: res.data.user
          }))
        }
      } catch (e) {
        console.log(e)
      }
    })
    setUserList(users)
  }

  const contextValue: HttpContextProps = {
    busList,
    loadBusList,
    signUp,
    userDatas,
    userList,
    setUserList: setUserListFn
  }
  return (
    <HttpContext.Provider value={contextValue}>{children}</HttpContext.Provider>
  )
}
