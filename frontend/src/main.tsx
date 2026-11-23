import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Navigate, Routes, Route } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./api/client";
import "./index.css";
import App from "./App";
import Greeting from "./greeting/greeting";
import Extracurriculars from "./extracurricular/Extracurriculars";
import CreateExtracurricular from "./extracurricular/CreateExtracurricular";
import ViewExtracurricular from "./extracurricular/ViewExtracurriculars";
import EditExtracurricular from "./extracurricular/EditExtracurricular";
import Essay from "./essays/essays";
import CounselorCollegeList from "./college/counselor/college-list.tsx";
import CounselorCollegeDetails from "./college/counselor/college-details.tsx";
import CounselorAddCollege from "./college/counselor/add-college.tsx";
import CounselorStudentCollegeList from "./college/counselor/student-college-list.tsx";
import MyCollegeList from "./college/student/my-college-list";
import SchoolDetails from "./college/student/school-details";
import UploadVideo from "./video/upload-video";
import VideoList from "./video/videos";
import Student from "./student/student";
import Counselor from "./counselor/counselor";
import CounselorStudents from "./counselor/counselor-students";
import PendingReviews from "./counselor/pending-reviews";
import Layout from "./layout/layout";
import StudentMediaList from "./media/student-media-list";
import CounselorMediaList from "./media/counselor/media-list";
import CounselorMediaDetails from "./media/counselor/media-details";
import CounselorAddMedia from "./media/counselor/add-media";
import UpdateNotificationPreferences from "./notifications/update-notification-preferences";
import CreateUser from "./user/create-user";
import LogOut from "./logout/logout";
import LogIn from "./login/login";
import SearchStudents from "./counselor/search-students.tsx";
import { AuthProvider } from "./contexts/auth-context.tsx";
import {
  ACCOUNT_EXPIRED_PATH,
  CLAIM_ACCOUNT_PATH,
} from "./contexts/auth-context.ts";
import AccountExpired from "./account-expiration/account-expired";
import ClaimAccount from "./claim-account/claim-account";
import StudentTaskList from "./task/student-task-list.tsx";
import StudentProfile from "./student/student-profile.tsx";
import StudentEssays from "./counselor/student-essays.tsx";
import AssignTaskForm from "./task/assign-task-form.tsx";
import EditTaskForm from "./task/edit-task-form.tsx";
import CounselorStudentTasks from "./counselor/counselor-student-tasks.tsx";
import EditProfile from "./profile/edit-profile.tsx";
import EditProfilePicture from "./profile/edit-profile-picture.tsx";
import ViewProfile from "./profile/view-profile.tsx";
import Chats from "./chat/chats.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path={ACCOUNT_EXPIRED_PATH} element={<AccountExpired />} />
            <Route path={CLAIM_ACCOUNT_PATH} element={<ClaimAccount />} />
            <Route path="/login" element={<LogIn />} />
            <Route
              path="/student/login"
              element={<LogIn audience="student" />}
            />
            <Route
              path="/counselor/login"
              element={<LogIn audience="counselor" />}
            />
            <Route element={<Layout />}>
              <Route path="/" element={<App />} />
              <Route path="/greeting" element={<Greeting />} />
              <Route path="/extracurriculars" element={<Extracurriculars />} />
              <Route
                path="/extracurriculars/new"
                element={<CreateExtracurricular />}
              />
              <Route
                path="/extracurriculars/:id"
                element={<ViewExtracurricular />}
              />
              <Route
                path="/extracurriculars/:id/edit"
                element={<EditExtracurricular />}
              />
              <Route path="/essay" element={<Essay />} />
              <Route
                path="/counselor-college-list"
                element={<CounselorCollegeList />}
              />
              <Route
                path="/counselor-college-list/new"
                element={<CounselorAddCollege />}
              />
              <Route
                path="/counselor-college-list/:id"
                element={<CounselorCollegeDetails />}
              />
              {/* old student college search page; Add School on My List replaced it */}
              <Route
                path="/student-college-list/*"
                element={<Navigate to="/college-applications" replace />}
              />
              <Route path="/college-applications" element={<MyCollegeList />} />
              <Route
                path="/college-applications/:id"
                element={<SchoolDetails />}
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
              {/* old upload page; Create Module replaced it */}
              <Route
                path="/upload-media"
                element={<Navigate to="/counselor-media-list/new" replace />}
              />
              <Route
                path="/student-media-list"
                element={<StudentMediaList />}
              />
              <Route
                path="/counselor-media-list"
                element={<CounselorMediaList />}
              />
              <Route
                path="/counselor-media-list/new"
                element={<CounselorAddMedia />}
              />
              <Route
                path="/counselor-media-list/:id"
                element={<CounselorMediaDetails />}
              />
              <Route
                path="/update-notification-preferences"
                element={<UpdateNotificationPreferences />}
              />
              <Route path="/create-account" element={<CreateUser />} />
              <Route path="/logout" element={<LogOut />} />
              <Route path="/student-tasks" element={<StudentTaskList />} />
              <Route path="/students/:id" element={<StudentProfile />} />
              <Route
                path="/counselor/students/:studentId/essays"
                element={<StudentEssays />}
              />
              <Route
                path="/counselor/students/:studentId/colleges"
                element={<CounselorStudentCollegeList />}
              />
              <Route path="/search" element={<SearchStudents />} />
              <Route path="/assign-task" element={<AssignTaskForm />} />
              <Route
                path="/counselor/students/:id"
                element={<CounselorStudentTasks />}
              />
              <Route path="/profile" element={<EditProfile />} />
              <Route
                path="/counselor/students/:studentId/tasks/:taskId/edit"
                element={<EditTaskForm />}
              />
              <Route path="/profile/picture" element={<EditProfilePicture />} />
              <Route path="/profile/:userId" element={<ViewProfile />} />
              <Route path="/chats" element={<Chats />} />
              <Route path="/chats/:userId" element={<Chats />} />
            </Route>
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
);
