import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useHttp } from '../Context/HttpContext'
import { useAuth } from '../Context/AuthContext'

const Home = () => {
  const { userData, session } = useAuth()
  const { loadBusList, busList } = useHttp()
  const [search, setSearch] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)
  const Navigate = useNavigate()
  useEffect(() => {
    setLoading(true)
    loadBusList(search).finally(() => {
      setLoading(false)
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search, session])

  return (
    <div className="flex flex-col w-screen overflow-hidden items-center">
      <div className="max-w-[90rem] w-[80vw] mt-12 flex flex-col gap-3">
        <form className="flex items-center gap-2 justify-between">
          <input
            type="text"
            placeholder="Search"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="bg-gray-100 px-2 h-8 outline-none"
          />
          {userData?.busOwner && (
            <button
              className="px-2 h-8 bg-gray-500 text-white"
              onClick={(e) => {
                e.preventDefault()
                Navigate('/bus')
              }}
            >
              Admin
            </button>
          )}
        </form>
        <div>
          {loading ? (
            <div className="spinner"></div>
          ) : (
            <ul className="flex flex-wrap">
              {busList.length > 0 ? (
                busList.map((b) => {
                  return (
                    <li
                      key={`${b.ID + '-homePage-List'}`}
                      onClick={() => {
                        Navigate(`/user/${b.ID}`)
                      }}
                      className="w-64 h-32 bg-gray-100 rounded p-2"
                    >
                      <h3>{b.name}</h3>
                      <p>Lat: {b.lat}</p>
                      <p>Long: {b.long}</p>
                    </li>
                  )
                })
              ) : (
                <div>No Bus Found</div>
              )}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}

export default Home
