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

const Bus = () => {
  const [map, setMap] = useState<Map | null>(null)
  const [oldMarker, setOldMarker] = useState<Feature | null>(null)
  const [zoom, setZoom] = useState(12)
  const [load, setLoad] = useState(false)
  const { connected, customAlert, setUserSocket, location } = useWebSocket()

  const params = useParams()

  useEffect(() => {
    if (!load) {
      setUserSocket(params.busId || '1')
      setLoad(true)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.busId])

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
      setZoom(newZoom || 12)
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
    <div className="flex flex-col justify-center items-center">
      <p>{connected ? 'Connected' : 'Disconnected'}</p>
      <p>{customAlert}</p>
      <p>{'Latitude: ' + location.lat}</p>
      <p>{'Longitude: ' + location.long}</p>
      <div
        id="user-map"
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
