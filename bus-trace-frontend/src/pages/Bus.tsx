import { useEffect, useState } from 'react'
import { useWebSocket } from '../Context/WebContext'
import { useParams } from 'react-router-dom'
import Map from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import XYZ from 'ol/source/XYZ'
import { fromLonLat } from 'ol/proj'
import { useAuth } from '../Context/AuthContext'
import Alert from '../components/Alert'
import { useHttp } from '../Context/HttpContext'
import { HiMenuAlt2 } from 'react-icons/hi'
import { FaChevronLeft } from 'react-icons/fa'

const Bus = () => {
  const { session } = useAuth()
  const { userList, userDatas } = useHttp()
  const { customAlert, connected, setBusSocket, busClose } = useWebSocket()
  const [map, setMap] = useState<Map | null>(null)
  const [zoom, setZoom] = useState(16)
  const [load, setLoad] = useState(false)
  const [isMenuOpened, setIsMenuOpened] = useState(false)
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

    map?.getView().on('change:resolution', () => {
      const newZoom = map?.getView().getZoom()
      setZoom(newZoom!)
    })
    setMap(map)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    const onLoad = async () => {
      try {
        const position: GeolocationPosition = await new Promise(function (
          resolve,
          reject
        ) {
          navigator.geolocation.getCurrentPosition(
            (position) => resolve(position),
            (error) => reject(error)
          )
        })

        map?.setView(
          new View({
            center: fromLonLat([
              position.coords.longitude,
              position.coords.latitude
            ]),
            zoom: zoom
          })
        )
      } catch (e) {
        console.error(e)
      }
    }
    onLoad()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [map, zoom])

  return (
    <div className="flex flex-col justify-center items-center h-screen w-screen relative">
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
      <div className={`fixed top-0 left-0 z-20 `}>
        <div
          className={`bg-white max-w-[30rem] w-[90vw] h-screen top-0 fixed  ${isMenuOpened ? ' left-0' : '-left-[100%]'} transition-all duration-300 z-30`}
        >
          <ul className="grid overflow-auto">
            {(userList.length > 0 ? userList : []).map((userId) => {
              const user = Object.keys(userDatas).find((key) => key === userId)
              if (user) {
                const userData = userDatas[user]
                return (
                  <>
                    <li
                      key={userId + 'user'}
                      className="bg-white shadow-md rounded-lg p-4 mb-4"
                    >
                      <div>
                        <strong>Name:</strong> {userData.name}
                      </div>
                      <div>
                        <strong>Phone:</strong> {userData.phone}
                      </div>
                      <div>
                        <strong>Email:</strong> {userData.email}
                      </div>
                      <div>
                        <strong>Location:</strong> Latitude: {userData.lat},
                        Longitude: {userData.long}
                      </div>
                      <div>Location: {userData.location}</div>
                    </li>
                  </>
                )
              } else {
                return null
              }
            })}
          </ul>
          <FaChevronLeft
            className="absolute top-0 -right-7 h-8 w-8 p-1 bg-gray-200"
            onClick={() => setIsMenuOpened(false)}
          />
        </div>
        <HiMenuAlt2
          className="bg-white h-8 w-8"
          onClick={() => setIsMenuOpened(true)}
        />
      </div>
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
