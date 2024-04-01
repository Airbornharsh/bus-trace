import { useState } from 'react'
import { supabase } from '../helpers/Supabase'
import { useHttp } from '../Context/HttpContext'
import { useNavigate } from 'react-router-dom'

const SignUp = () => {
  const { signUp } = useHttp()
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    password: '',
    confirmPassword: ''
  })
  const Navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent<HTMLButtonElement>) => {
    e.preventDefault()
    try {
      const { data, error } = await supabase.auth.signUp({
        email: formData.email,
        password: formData.password,
        phone: formData.phone,
        options: {
          emailRedirectTo: 'https://bus-trace.harshkeshri.com/login',
        }
      })
      if (error) {
        console.error('Error:', error.message)
        alert(error.message)
        return
      }
      if (data.user) {
        console.log('Data:', data)
        await signUp(
          data.user.id,
          formData.name,
          formData.email,
          formData.phone
        )
      }
      console.log('Data:', data)
      Navigate('/login')
    } catch (e) {
      console.error('Error:', e)
    }
  }

  return (
    <div className="flex justify-center items-center w-screen h-screen bg-gray-200">
      <div className="w-[90vw] max-w-[18rem] bg-white min-h-[20rem] rounded shadow-lg">
        <form className="p-2 flex flex-col py-4 gap-3">
          <div>
            <label htmlFor="name" className="text-sm">
              Name
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded h-8"
            />
          </div>
          <div>
            <label htmlFor="phone" className="text-sm">
              Phone
            </label>
            <input
              type="text"
              id="phone"
              name="phone"
              value={formData.phone}
              onChange={(e) =>
                setFormData({ ...formData, phone: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded h-8"
            />
          </div>
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
          <div>
            <label htmlFor="confirmPassword" className="text-sm">
              Confirm Password
            </label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={(e) =>
                setFormData({ ...formData, confirmPassword: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded h-8"
            />
          </div>
          <div className="flex flex-col gap-1">
            <p
              className="text-xs text-blue-500 cursor-pointer"
              onClick={() => {
                Navigate('/login')
              }}
            >
              I have an Account?
            </p>
            <button
              type="submit"
              className="w-full p-1 bg-blue-500 text-white rounded"
              onClick={handleSubmit}
            >
              Sign Up
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default SignUp
