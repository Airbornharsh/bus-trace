import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Home from './pages/Home'
import User from './pages/User'
import Bus from './pages/Bus'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/user/:userId/:busId" element={<User />} />
        <Route path="/bus/:userId/:busId" element={<Bus />} />
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
