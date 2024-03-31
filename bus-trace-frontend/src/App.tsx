import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Home from './pages/Home'
import User from './pages/User'
import Bus from './pages/Bus'
import SignUp from './pages/SignUp'
import Login from './pages/Login'
import { useEffect } from 'react'
import { supabase } from './helpers/Supabase'

const App = () => {
  useEffect(() => {
    const onload = async () => {
      const { data, error } = await supabase.auth.getUser()
      if (error) {
        console.error('Error:', error.message)
        window.location.href = '/login'
        return
      }
      if (!data) {
        window.location.href = '/login'
      }
    }
    if (
      window.location.pathname !== '/login' &&
      window.location.pathname !== '/signup'
    ) {
      onload()
    }
  }, [])

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/login" element={<Login />} />
        <Route path="/user/:busId" element={<User />} />
        <Route path="/bus" element={<Bus />} />
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
