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

const User = () => {
  const [map, setMap] = useState<Map | null>(null)
  const [oldMarker, setOldMarker] = useState<Feature | null>(null)
  const [zoom, setZoom] = useState(16)
  const [load, setLoad] = useState(false)
  const { session } = useAuth()
  const { connected, customAlert, setUserSocket, location } = useWebSocket()

  const params = useParams()

  useEffect(() => {
    if (!load && params.busId && session?.user) {
      setUserSocket(params.busId)
      setLoad(true)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.busId, session])

  useEffect(() => {
    const map = new Map({
      target: 'user-map',
      layers: [
        new TileLayer({
          source: new XYZ({
            url: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
          })
        })
      ]
    })

    const marker = new Feature({
      geometry: new Point(fromLonLat([location.long, location.lat]))
    })

    // Custom marker style
    // marker.setStyle(

    // )

    // Adding marker to a vector layer
    const vectorLayer = new VectorLayer({
      source: new VectorSource({
        features: [marker]
      })
    })

    map?.addLayer(vectorLayer)

    map?.getView().on('change:resolution', () => {
      const newZoom = map?.getView().getZoom()
      setZoom(newZoom!)
    })

    setOldMarker(marker)

    setMap(map)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    map?.setView(
      new View({
        center: fromLonLat([location.long, location.lat]),
        zoom: zoom
      })
    )
    oldMarker?.setGeometry(new Point(fromLonLat([location.long, location.lat])))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location, map, zoom])

  return (
    <div className="flex flex-col justify-center items-center h-screen w-screen relative">
      <div className="fixed top-0 -translate-x-[50%] left-[50%] z-20">
        {customAlert && <Alert customAlert={customAlert} />}
      </div>
      <p className="fixed top-0 right-0 p-1 z-20 bg-white">
        {connected ? (
          <img src="/icon/connected.png" alt="connected" className="w-6 h-6" />
        ) : (
          <img
            src="/icon/disconnected.png"
            alt="disconnected"
            className="w-6 h-6"
          />
        )}
      </p>
      <div className="fixed bottom-0 -translate-x-[50%] left-[50%] z-20 text-xs">
        <p>{'Lat: ' + location.lat}</p>
        <p>{'Lon: ' + location.long}</p>
      </div>
      <div
        id="user-map"
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

export default User
