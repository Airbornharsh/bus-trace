import React, { createContext, useContext, useState, ReactNode } from 'react'
import { validate } from '../helpers/Json'
import { useAuth } from './AuthContext'

interface Position {
  busId: string
  lat: string
  long: string
}

interface Location {
  lat: number
  long: number
}

interface WebSocketContextProps {
  connected: boolean
  socket: WebSocket | null
  location: Location
  userLocations: { [key: string]: Location }
  userList: string[]
  customAlert: string
  setBusSocket: () => void
  setUserSocket: (busId: string) => void
  sendMessage: (message: Position) => void
}

const WebSocketContext = createContext<WebSocketContextProps | undefined>(
  undefined
)

// eslint-disable-next-line react-refresh/only-export-components
export const useWebSocket = () => {
  const context = useContext(WebSocketContext)

  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }

  return context
}

interface WebSocketProviderProps {
  children: ReactNode
}

const WebSocketUrl = 'ws://localhost:8000/ws'
// const WebSocketUrl = 'wss://bus-trace-websocket-server.onrender.com/ws'

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({
  children
}) => {
  const { session } = useAuth()
  const [connected, setConnected] = useState(false)
  const [socket, setSocket] = useState<WebSocket | null>(null)
  const [customAlert, setCustomAlert] = useState<string>('')
  const [location, setLocation] = useState<Location>({
    lat: 0,
    long: 0
  })
  const [userLocations, setUserLocations] = useState<{
    [key: string]: Location
  }>({})
  const [userList, setUserList] = useState<string[]>([])

  const setBusSocketFn = () => {
    console.log('socket Connection')
    if (socket) {
      socket.close()
    }
    const newSocket = new WebSocket(
      `${WebSocketUrl}/bus/${session?.access_token}`
    )

    newSocket.addEventListener('open', (event) => {
      setConnected(true)
      console.log('WebSocket connection opened:', event)
      const data = 0
      setInterval(() => {
        if (navigator.geolocation) {
          navigator.geolocation.getCurrentPosition(
            (position) => {
              newSocket.send(
                JSON.stringify({
                  busId: '1',
                  lat: position.coords.latitude + data,
                  long: position.coords.longitude + data
                })
              )
            },
            (e) => {
              console.error(e)
            },
            { enableHighAccuracy: true }
          )
        } else {
          console.log('Not Supported by Browser')
          alert('Not supported')
        }
      }, 5000)
    })

    newSocket.addEventListener('close', (event) => {
      setConnected(false)
      console.log('WebSocket connection closed:', event)
    })

    newSocket.addEventListener('error', (event) => {
      setConnected(false)
      console.error('WebSocket connection error:', event)
    })

    newSocket.addEventListener('message', (event) => {
      console.log('WebSocket message:', event)
      if (validate(event.data)) {
        if (
          event.data.includes('lat') &&
          event.data.includes('long') &&
          event.data.includes('userId')
        ) {
          const data = JSON.parse(event.data)
          setUserLocations((prev) => ({
            ...prev,
            [data.userId]: {
              lat: data.lat,
              long: data.long
            }
          }))
        } else {
          const data = JSON.parse(event.data)
          setUserList(data || [])
        }
      } else if (event.data.split(' ')[0] === 'Status:') {
        setCustomAlert(event.data)
        setTimeout(() => {
          setCustomAlert('')
        }, 5000)
      }
    })

    setSocket(newSocket)
  }

  const setUserSocketFn = (busId: string) => {
    console.log('socket Connection')
    if (socket) {
      socket.close()
    }

    const newSocket = new WebSocket(
      `${WebSocketUrl}/user/${busId}/${session?.access_token}`
    )

    newSocket.addEventListener('open', (event) => {
      setConnected(true)
      console.log('WebSocket connection opened:', event)
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            newSocket.send(
              JSON.stringify({
                lat: position.coords.latitude,
                long: position.coords.longitude
              })
            )
          },
          (e) => {
            console.error(e)
          },
          { enableHighAccuracy: true }
        )
      } else {
        console.log('Not Supported by Browser')
        alert('Not supported')
      }
      // setInterval(() => {
      //   if (navigator.geolocation) {
      //     navigator.geolocation.getCurrentPosition(
      //       (position) => {
      //         newSocket.send(
      //           JSON.stringify({
      //             lat: position.coords.latitude,
      //             long: position.coords.longitude
      //           })
      //         )
      //       },
      //       (e) => {
      //         console.error(e)
      //       },
      //       { enableHighAccuracy: true }
      //     )
      //   } else {
      //     console.log('Not Supported by Browser')
      //     alert('Not supported')
      //   }
      // }, 5000)
    })

    newSocket.addEventListener('close', (event) => {
      setConnected(false)
      console.log('WebSocket connection closed:', event)
    })

    newSocket.addEventListener('error', (event) => {
      setConnected(false)
      console.error('WebSocket connection error:', event)
    })

    newSocket.addEventListener('message', (event) => {
      console.log('WebSocket message:', event.data)
      if (event.data.includes('lat') && event.data.includes('long')) {
        const data = JSON.parse(event.data)
        setLocation({
          lat: data.lat,
          long: data.long
        })
      } else if (event.data.split(' ')[0] === 'Status:') {
        setCustomAlert(event.data)
        setTimeout(() => {
          setCustomAlert('')
        }, 5000)
      }
    })

    setSocket(newSocket)
  }

  const sendMessage = (message: Position) => {
    console.log(socket?.readyState)
    if (socket) {
      console.log(message)
      socket.send(JSON.stringify(message))
    }
  }

  const contextValue: WebSocketContextProps = {
    connected,
    socket,
    setBusSocket: setBusSocketFn,
    location,
    userLocations,
    userList,
    customAlert,
    sendMessage,
    setUserSocket: setUserSocketFn
  }

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  )
}
