import { useEffect, useState } from 'react'
import { useWebSocket } from '../Context/WebContext'
import { useParams } from 'react-router-dom'

const Bus = () => {
  const [load, setLoad] = useState(false)
  const { connected, setBusSocket } = useWebSocket()

  const params = useParams()

  useEffect(() => {
    if (!load) {
      setBusSocket(params.busId || '1')
      setLoad(true)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.busId])

  return (
    <div>
      <h2>Admin</h2>
      <p>{connected ? 'Connected' : 'Disconnected'}</p>
    </div>
  )
}

export default Bus
