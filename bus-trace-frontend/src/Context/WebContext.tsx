import React, { createContext, useContext, useState, ReactNode } from 'react'

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
  setBusSocket: (id: string) => void
  setUserSocket: (id: string) => void
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
  const [connected, setConnected] = useState(false)
  const [socket, setSocket] = useState<WebSocket | null>(null)
  const [location, setLocation] = useState<Location>({
    lat: 0,
    long: 0
  })

  const setBusSocketFn = (id: string) => {
    console.log('socket Connection')
    if (socket) {
      socket.close()
    }

    const newSocket = new WebSocket(`${WebSocketUrl}/bus/1/${id}`)

    newSocket.addEventListener('open', (event) => {
      setConnected(true)
      console.log('WebSocket connection opened:', event)
      const data = 0
      setInterval(() => {
        if (navigator.geolocation) {
          navigator.geolocation.getCurrentPosition(
            (position) => {
              console.log('Latitude:', position.coords.latitude)
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
    })

    setSocket(newSocket)
  }

  const setUserSocketFn = (id: string) => {
    console.log('socket Connection')
    if (socket) {
      socket.close()
    }

    const newSocket = new WebSocket(`${WebSocketUrl}/user/2/${id}`)

    newSocket.addEventListener('open', (event) => {
      setConnected(true)
      console.log('WebSocket connection opened:', event)
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
      }
    })

    setSocket(newSocket)
  }
  // useEffect(() => {
  //   setBusSocketFn('1')
  //   // eslint-disable-next-line react-hooks/exhaustive-deps
  // }, [])

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
    sendMessage,
    setUserSocket: setUserSocketFn
  }

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  )
}
