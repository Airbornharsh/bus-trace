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
      const { data, error } = await supabase.auth.signInWithPassword({
        email: formData.email,
        password: formData.password
      })
      if (error) {
        console.error('Error:', error.message)
        alert(error.message)
        return
      }
      console.log('Response:', data)
      Navigate('/')
    } catch (error) {
      console.error('Error:', error)
    }
  }

  return (
    <div className="flex justify-center items-center w-screen h-screen bg-gray-200">
      <div className="w-[90vw] max-w-[18rem] bg-white  rounded shadow-lg">
        <form className="p-2 flex flex-col py-4 gap-3">
          <div>
            <label htmlFor="email" className="text-sm">
              Email
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={(e) =>
                setFormData({ ...formData, email: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded h-8"
            />
          </div>
          <div>
            <label htmlFor="password" className="text-sm">
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={(e) =>
                setFormData({ ...formData, password: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded h-8"
            />
          </div>
          <div className="flex flex-col gap-1">
            <p
              className="text-xs text-blue-500 cursor-pointer"
              onClick={() => {
                Navigate('/signup')
              }}
            >
              Create a New Account
            </p>
            <button
              type="submit"
              className="w-full p-1 bg-blue-500 text-white rounded"
              onClick={handleSubmit}
            >
              Login
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default SignUp
