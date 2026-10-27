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
import CounselorCollegeList from "./college/counselor/college-list.tsx";
import CounselorCollegeDetails from "./college/counselor/college-details.tsx";
import StudentCollegeList from "./college/student/college-list.tsx";
import StudentCollegeDetails from "./college/student/college-details.tsx";
import MyCollegeList from "./college/student/my-college-list";
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
import { ACCOUNT_EXPIRED_PATH } from "./contexts/auth-context.ts";
import AccountExpired from "./account-expiration/account-expired";
import StudentProfile from "./student/student-profile.tsx";

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
              <Route
                path="/counselor-college-list"
                element={<CounselorCollegeList />}
              />
              <Route
                path="/counselor-college-list/:id"
                element={<CounselorCollegeDetails />}
              />
              <Route
                path="/student-college-list"
                element={<StudentCollegeList />}
              />
              <Route
                path="/student-college-list/:id"
                element={<StudentCollegeDetails />}
              />
              <Route path="/college-applications" element={<MyCollegeList />} />
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
              <Route path="/students/:id" element={<StudentProfile />} />
            </Route>
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
);
