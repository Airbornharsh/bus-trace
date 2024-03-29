import { useState } from 'react'
import { supabase } from '../helpers/Supabase'
import { useNavigate } from 'react-router-dom'

const SignUp = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  })
  const Navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent<HTMLButtonElement>) => {
    e.preventDefault()
    try {
      const res = await supabase.auth.signInWithPassword({
        email: formData.email,
        password: formData.password
      })

      console.log('Response:', res)
      Navigate('/')
    } catch (error) {
      console.error('Error:', error)
    }
  }

  return (
    <div className="flex justify-center items-center w-screen h-screen bg-gray-200">
      <div className="w-[90vw] max-w-[25rem] bg-white  rounded shadow-lg">
        <form className="p-2 flex-col">
          <div>
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={(e) =>
                setFormData({ ...formData, email: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded"
            />
          </div>
          <div>
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={(e) =>
                setFormData({ ...formData, password: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded"
            />
          </div>
          <button
            type="submit"
            className="w-full p-1 bg-blue-500 text-white rounded"
            onClick={handleSubmit}
          >
            Login
          </button>
        </form>
      </div>
    </div>
  )
}

export default SignUp
