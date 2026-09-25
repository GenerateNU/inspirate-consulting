import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import App from './App'
import Greeting from './greeting/greeting'
import CreateCollege from './college/create-college'
import CollegeList from './college/college-list'
import Student from './student/student'
import Counselor from './counselor/counselor'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/greeting" element={<Greeting />} />
        <Route path="/create-college" element={<CreateCollege />} />
        <Route path="/college-list" element={<CollegeList />} />
        <Route path="/student" element={<Student />} />
        <Route path="/counselor" element={<Counselor />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)