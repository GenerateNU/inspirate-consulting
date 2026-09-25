import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import App from './App'
import Greeting from './greeting/greeting'
import Essay from "./essays/essays";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/greeting" element={<Greeting />} />
        <Route path="/essay" element={<Essay />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);