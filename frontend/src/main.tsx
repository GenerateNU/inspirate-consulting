import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import App from './App'
import Greeting from './greeting/greeting'
import CreateTask from './task/CreateTask'
import TaskList from './task/TaskList'
import CreateCollege from './college/create-college'
import CollegeList from './college/college-list'
import CreateCollegeApplication from './college/create-application'
import CollegeApplicationList from './college/colleges'
import VideoList from './video/videos'
import UploadVideo from './video/upload-video'



createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/greeting" element={<Greeting />} />
        <Route path="/create-task" element={<CreateTask />} />
        <Route path="/tasks" element={<TaskList />} />
        <Route path="/create-college" element={<CreateCollege />} />
        <Route path="/college-list" element={<CollegeList />} />
        <Route path="/create-application" element={<CreateCollegeApplication />} />
        <Route path="/colleges" element={<CollegeApplicationList />} />
        <Route path="/upload-video" element={<UploadVideo />} />
        <Route path="/videos" element={<VideoList />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)