import { useEffect, useState } from 'react'
import { useAuth } from '../Context/AuthContext'
import axios from 'axios'
import { useNavigate } from 'react-router-dom'

const Home = () => {
  const { session } = useAuth()
  const [buses, setBuses] = useState<
    {
      ID: string
      name: string
      lat: number
      long: number
    }[]
  >([])
  const [search, setSearch] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)
  const Navigate = useNavigate()
  useEffect(() => {
    const onLoad = async () => {
      setLoading(true)
      try {
        const res = await axios.get(`http://localhost:8001/bus?search=a`)
        setBuses(res.data.buses)
      } catch (e) {
        console.log(e)
      } finally {
        setLoading(false)
      }
    }

    if (session?.user) {
      onLoad()
    }
  }, [session, search])

  const searchHandler = async (e: React.FormEvent<HTMLButtonElement>) => {
    e.preventDefault()
    setLoading(true)
    try {
      const res = await axios.get(`http://localhost:8001/bus?search=${search}`)
      setBuses(res.data.buses)
    } catch (e) {
      console.log(e)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col">
      <form>
        <input
          type="text"
          placeholder="Search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        {!loading && (
          <button type="submit" onClick={searchHandler}>
            Search
          </button>
        )}
      </form>
      <div>
        <ul>
          {buses.map((b) => {
            return (
              <li
                key={`${b.ID + '-homePage-List'}`}
                onClick={() => {
                  Navigate(`/user/${b.ID}`)
                }}
              >
                {b.name}
              </li>
            )
          })}
        </ul>
      </div>
    </div>
  )
}

export default Home
