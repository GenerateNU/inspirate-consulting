import { AuthContext } from "./auth-context";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  useLoginUser,
  useLogoutUser,
  useResetUserPassword,
  type loginUserResponse,
} from "../api/endpoints";
import type { ErrorModel } from "../api/models";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [resetTime, setResetTime] = useState<string | null>(null);
  const [userID, setUserID] = useState<string | null>(null);
  const [counselorID, setCounselorID] = useState<string | null>(null);
  const [error, setError] = useState<Promise<ErrorModel> | null>(null);
  const [path, setPath] = useState("/");
  const [studentID, setStudentID] = useState<string | null>(null);
  const navigate = useNavigate();

  const { trigger: loginFunc } = useLoginUser();
  const { trigger: logoutFunc } = useLogoutUser();
  const { trigger: passFunc } = useResetUserPassword();

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
          void navigate(path);
        },
        onError: (err) => {
          setError(err);
          console.log("login did not work", error);
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
          void navigate(path);
          setPath("/");
        },

        onError: (err) => {
          setError(err);
        },
      },
    );
  };

  const new_password = (new_password: string) => {
    return passFunc(
      { new_password },
      {
        onSuccess: () => {
          void navigate(path);
        },
        onError: (err) => {
          setError(err);
          console.log("resetting password did not work", error);
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
