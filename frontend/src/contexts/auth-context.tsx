import { AuthContext } from "./auth-context";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  useFetchCurrentUser,
  useLoginUser,
  useLogoutUser,
  useResetUserPassword,
  type loginUserResponse,
} from "../api/endpoints";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [resetTime, setResetTime] = useState<string | null>(null);
  const [userID, setUserID] = useState<string | null>(null);
  const [counselorID, setCounselorID] = useState<string | null>(null);
  const [studentID, setStudentID] = useState<string | null>(null);
  const navigate = useNavigate();

  const { trigger: loginFunc } = useLoginUser();
  const { trigger: logoutFunc } = useLogoutUser();
  const { trigger: passFunc } = useResetUserPassword();

  const { data: currentUser, mutate: mutateCurrentUser } = useFetchCurrentUser({
    swr: { shouldRetryOnError: false, revalidateOnFocus: false },
  });

  useEffect(() => {
    if (currentUser?.status !== 200) {
      return;
    }
    setResetTime(currentUser.data.reset_time);
    setUserID(currentUser.data.id);
    setCounselorID(currentUser.data.counselor_id);
    setStudentID(currentUser.data.student_id);
  }, [currentUser]);

  const login = (email: string, password: string) => {
    return loginFunc(
      { email, password },
      {
        onSuccess: (resp: loginUserResponse) => {
          if (resp.status != 200) {
            return;
          }
          setResetTime(resp.data.reset_time);
          setUserID(resp.data.id);
          setCounselorID(resp.data.counselor_id);
          setStudentID(resp.data.student_id);
          console.log(resp.data.id);
          void navigate("/");
        },
        onError: () => {
          console.log("login did not work");
          return;
        },
      },
    );
  };

  const logout = () => {
    return logoutFunc(
      {},
      {
        onSuccess: () => {
          setResetTime(null);
          setUserID(null);
          setCounselorID(null);
          setStudentID(null);
          void mutateCurrentUser(undefined, { revalidate: false });
          void navigate("/");
        },
      },
    );
  };

  const new_password = (new_password: string) => {
    return passFunc(
      { new_password },
      {
        onSuccess: () => {
          void navigate("/");
        },
        onError: () => {
          console.log("resetting password did not work");
        },
      },
    );
  };

  return (
    <AuthContext
      value={{
        resetTime: resetTime,
        userID: userID,
        counselorID: counselorID,
        studentID: studentID,
        login: login,
        logout: logout,
        reset_password: new_password,
      }}
    >
      {children}
    </AuthContext>
  );
}
