import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import './index.css'
import App from './App'
import Greeting from './greeting/greeting'
import Essay from "./essays/essays";
import CreateCollege from './college/create-college'
import CollegeList from './college/college-list'
import CreateCollegeApplication from './college/create-application'
import CollegeApplicationList from './college/colleges'
import UploadVideo from "./video/upload-video";
import VideoList from "./video/videos";
import Student from "./student/student";
import Counselor from "./counselor/counselor";
import CreateTask from "./task/CreateTask";
import TaskList from "./task/TaskList";
import Layout from "./layout/layout";


createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<App />} />
          <Route path="/greeting" element={<Greeting />} />
          <Route path="/essay" element={<Essay />} />
          <Route path="/create-college" element={<CreateCollege />} />
          <Route path="/college-list" element={<CollegeList />} />
          <Route
            path="/create-application"
            element={<CreateCollegeApplication />}
          />
          <Route path="/colleges" element={<CollegeApplicationList />} />
          <Route path="/upload-video" element={<UploadVideo />} />
          <Route path="/videos" element={<VideoList />} />
          <Route path="/student" element={<Student />} />
          <Route path="/counselor" element={<Counselor />} />
          <Route path="/create-task" element={<CreateTask />} />
          <Route path="/tasks" element={<TaskList />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);
