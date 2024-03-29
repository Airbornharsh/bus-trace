import { useEffect, useState } from 'react'
import { useWebSocket } from '../Context/WebContext'
import { useParams } from 'react-router-dom'
import Map from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import XYZ from 'ol/source/XYZ'
import Feature from 'ol/Feature'
import Point from 'ol/geom/Point'
import { fromLonLat } from 'ol/proj'
import { Vector as VectorLayer } from 'ol/layer'
import { Vector as VectorSource } from 'ol/source'
import { useAuth } from '../Context/AuthContext'

const Bus = () => {
  const [map, setMap] = useState<Map | null>(null)
  const [zoom, setZoom] = useState(16)
  const [load, setLoad] = useState(false)
  const { session } = useAuth()
  const {
    customAlert,
    userLocations,
    userList,
    connected,
    setBusSocket,
    busClose
  } = useWebSocket()

  const params = useParams()

  useEffect(() => {
    if (!load && session?.user) {
      setBusSocket()
      setLoad(true)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.busId, session])

  useEffect(() => {
    const map = new Map({
      target: 'bus-map',
      layers: [
        new TileLayer({
          source: new XYZ({
            url: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
          })
        })
      ]
    })

    const tempMarkers: Feature<Point>[] = []
    Object.keys(userLocations).forEach(() => {
      const marker = new Feature({
        geometry: new Point(fromLonLat([84.8535844, 22.260423]))
      })
      tempMarkers.push(marker)
    })
    const vectorLayer = new VectorLayer({
      source: new VectorSource({
        features: [...tempMarkers]
      })
    })

    map?.addLayer(vectorLayer)

    map?.getView().on('change:resolution', () => {
      const newZoom = map?.getView().getZoom()
      setZoom(newZoom || 12)
    })

    setMap(map)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userLocations])

  useEffect(() => {
    map?.setView(
      new View({
        center: fromLonLat([84.8535844, 22.260423]),
        zoom: zoom
      })
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location, map, zoom])

  return (
    <div className="flex flex-col justify-center items-center">
      <h2>Admin</h2>
      {connected ? (
        <button
          className="px-2 h-8 bg-gray-500 text-white"
          onClick={() => {
            busClose()
          }}
        >
          Close
        </button>
      ) : (
        <button
          className="px-2 h-8 bg-gray-500 text-white"
          onClick={() => {
            setBusSocket()
          }}
        >
          Reconnect
        </button>
      )}
      <p>{connected ? 'Connected' : 'Disconnected'}</p>
      <p>{customAlert}</p>
      <ul>
        {(userList.length > 0 ? userList : []).map((user) => (
          <li key={user}>{user}</li>
        ))}
      </ul>
      <div
        id="bus-map"
        style={{
          width: '90vw',
          height: '50rem',
          overflow: 'hidden'
        }}
      ></div>
    </div>
  )
}

export default Bus
