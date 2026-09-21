import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import App from './App'
import Greeting from './greeting/greeting'
import CreateExtracurricular from './extracurricular/CreateExtracurricular'
import Extracurriculars from './extracurricular/Extracurriculars'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/greeting" element={<Greeting />} />
        <Route path="/create-extracurricular" element={<CreateExtracurricular />} />
        <Route path="/extracurriculars" element={<Extracurriculars />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)