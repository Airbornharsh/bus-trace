import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Home from './pages/Home'
import User from './pages/User'
import Bus from './pages/Bus'
import SignUp from './pages/SignUp'
import Login from './pages/Login'
import { useEffect } from 'react'
import { createClient } from '@supabase/supabase-js'

const supabase = createClient(
  import.meta.env.VITE_APP_PROJECT_LINK,
  import.meta.env.VITE_APP_ANON_KEY
)

const App = () => {
  useEffect(() => {
    supabase.auth.onAuthStateChange((event, session) => {
      console.log('Event:', event)
      console.log('Session:', session)
    })
  }, [])
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/login" element={<Login />} />
        <Route path="/user/:busId" element={<User />} />
        <Route path="/bus/:busId" element={<Bus />} />
        <Route
          path="/*"
          element={
            <h1 className="flex justify-center items-center h-screen">
              Not Found
            </h1>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
