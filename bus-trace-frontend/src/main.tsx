import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { WebSocketProvider } from './Context/WebContext.tsx'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <WebSocketProvider>
    <App />
  </WebSocketProvider>
)
