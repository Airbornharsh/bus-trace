import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { WebSocketProvider } from './Context/WebContext.tsx'
import { AuthProvider } from './Context/AuthContext.tsx'
import { HttpProvider } from './Context/HttpContext.tsx'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <AuthProvider>
    <HttpProvider>
      <WebSocketProvider>
        <App />
      </WebSocketProvider>
    </HttpProvider>
  </AuthProvider>
)
