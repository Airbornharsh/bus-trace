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
import Alert from '../components/Alert'
import { useHttp } from '../Context/HttpContext'

const Bus = () => {
  const [map, setMap] = useState<Map | null>(null)
  const [zoom, setZoom] = useState(16)
  const [load, setLoad] = useState(false)
  const { session } = useAuth()
  const { customAlert, userLocations, connected, setBusSocket, busClose } =
    useWebSocket()
  const { userList, userDatas } = useHttp()

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
      setZoom(newZoom!)
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
    <div className="flex flex-col justify-center items-center h-screen w-screen relative">
      <div className="left top-0 right-0 z-20">
        <ul>
          {(userList.length > 0 ? userList : []).map((userId) => {
            const user = Object.keys(userDatas).find((key) => key === userId)
            if (user) {
              return (
                <li key={userId}>
                  <p>{userDatas[userId].name}</p>
                  <p>Lat: {userDatas[userId].lat || 0}</p>
                  <p>Long: {userDatas[userId].long || 0}</p>
                </li>
              )
            } else {
              return null
            }
          })}
        </ul>
      </div>
      <div className="fixed top-0 -translate-x-[50%] left-[50%] z-20">
        {customAlert && <Alert customAlert={customAlert} />}
      </div>
      <p className="fixed top-0 right-0 p-1 z-20 bg-white">
        {connected ? (
          <img
            src="/icon/connected.png"
            alt="connected"
            className="w-6 h-6"
            onClick={() => {
              busClose()
            }}
          />
        ) : (
          <img
            src="/icon/disconnected.png"
            alt="disconnected"
            className="w-6 h-6"
            onClick={() => {
              setBusSocket()
            }}
          />
        )}
      </p>
      <div
        id="bus-map"
        style={{
          width: '100vw',
          height: '100vh',
          overflow: 'hidden',
          zIndex: 0
        }}
      ></div>
    </div>
  )
}

export default Bus
