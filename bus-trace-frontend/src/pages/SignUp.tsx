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
        phone: formData.phone
        // options: {
        //   emailRedirectTo: 'http://localhost:3000'
        // }
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
      <div className="w-[90vw] max-w-[25rem] bg-white min-h-[20rem] rounded shadow-lg">
        <form className="p-2 flex-col">
          <div>
            <label htmlFor="name">Name</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded"
            />
          </div>
          <div>
            <label htmlFor="phone">Phone</label>
            <input
              type="text"
              id="phone"
              name="phone"
              value={formData.phone}
              onChange={(e) =>
                setFormData({ ...formData, phone: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded"
            />
          </div>
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
          <div>
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={(e) =>
                setFormData({ ...formData, confirmPassword: e.target.value })
              }
              className="w-full p-1 border border-gray-300 rounded"
            />
          </div>
          <button
            type="submit"
            className="w-full p-1 bg-blue-500 text-white rounded"
            onClick={handleSubmit}
          >
            {'Sign Up'}
          </button>
        </form>
      </div>
    </div>
  )
}

export default SignUp
