import { AuthContext } from "./auth-context";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  useLoginUser,
  useLogoutUser,
  useResetUserPassword,
} from "../api/endpoints";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [resetTime, setResetTime] = useState<string | null>(null);
  const [userID, setUserID] = useState<string | null>(null);
  const navigate = useNavigate();

  const { mutateAsync: loginFunc } = useLoginUser();
  const { mutateAsync: logoutFunc } = useLogoutUser();
  const { mutateAsync: passFunc } = useResetUserPassword();

  const login = (email: string, password: string) => {
    return loginFunc(
      { data: { email, password } },
      {
        onSuccess: (user) => {
          setResetTime(user.reset_time);
          setUserID(user.id);
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
    return logoutFunc(undefined, {
      onSuccess: () => {
        setResetTime(null);
        setUserID(null);
        void navigate("/");
      },
    });
  };

  const new_password = (new_password: string) => {
    return passFunc(
      { data: { new_password } },
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
        login: login,
        logout: logout,
        reset_password: new_password,
      }}
    >
      {children}
    </AuthContext>
  );
}
