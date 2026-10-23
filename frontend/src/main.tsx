import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./api/client";
import "./index.css";
import App from "./App";
import Greeting from "./greeting/greeting";
import CreateExtracurricular from "./extracurricular/CreateExtracurricular";
import Extracurriculars from "./extracurricular/Extracurriculars";
import Essay from "./essays/essays";
import CreateCollege from "./college/create-college";
import CollegeList from "./college/college-list";
import CreateCollegeApplication from "./college/create-application";
import UploadVideo from "./video/upload-video";
import VideoList from "./video/videos";
import Student from "./student/student";
import Counselor from "./counselor/counselor";
import CounselorStudents from "./counselor/counselor-students";
import PendingReviews from "./counselor/pending-reviews";
import CreateTask from "./task/CreateTask";
import TaskList from "./task/TaskList";
import Layout from "./layout/layout";
import UploadMedia from "./media/upload-media";
import StudentMediaList from "./media/student-media-list";
import CounselorMediaList from "./media/counselor-media-list";
import UpdateNotificationPreferences from "./notifications/update-notification-preferences";
import CreateUser from "./user/create-user";
import LogOut from "./logout/logout";
import LogIn from "./login/login";
import { AuthProvider } from "./contexts/auth-context.tsx";
import CollegeApplications from "./college/college-applications.tsx";
import { ACCOUNT_EXPIRED_PATH } from "./contexts/auth-context.ts";
import AccountExpired from "./account-expiration/account-expired";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path={ACCOUNT_EXPIRED_PATH} element={<AccountExpired />} />
            <Route element={<Layout />}>
              <Route path="/" element={<App />} />
              <Route path="/greeting" element={<Greeting />} />
              <Route
                path="/create-extracurricular"
                element={<CreateExtracurricular />}
              />
              <Route path="/extracurriculars" element={<Extracurriculars />} />
              <Route path="/essay" element={<Essay />} />
              <Route path="/create-college" element={<CreateCollege />} />
              <Route path="/college-list" element={<CollegeList />} />
              <Route
                path="/create-application"
                element={<CreateCollegeApplication />}
              />
              <Route
                path="/college-applications"
                element={<CollegeApplications />}
              />
              <Route path="/upload-video" element={<UploadVideo />} />
              <Route path="/videos" element={<VideoList />} />
              <Route path="/student" element={<Student />} />
              <Route path="/counselor" element={<Counselor />} />
              <Route
                path="/counselor-students"
                element={<CounselorStudents />}
              />
              <Route path="/pending-reviews" element={<PendingReviews />} />
              <Route path="/create-task" element={<CreateTask />} />
              <Route path="/tasks" element={<TaskList />} />
              <Route path="/upload-media" element={<UploadMedia />} />
              <Route
                path="/student-media-list"
                element={<StudentMediaList />}
              />
              <Route
                path="/counselor-media-list"
                element={<CounselorMediaList />}
              />
              <Route
                path="/update-notification-preferences"
                element={<UpdateNotificationPreferences />}
              />
              <Route path="/create-account" element={<CreateUser />} />
              <Route path="/login" element={<LogIn />} />
              <Route path="/logout" element={<LogOut />} />
            </Route>
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
);
