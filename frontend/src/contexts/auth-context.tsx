import { ACCOUNT_EXPIRED_PATH, AuthContext } from "./auth-context";
import { useEffect, useState } from "react";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import {
  useQueryClient,
  type MutationCacheNotifyEvent,
  type QueryCacheNotifyEvent,
} from "@tanstack/react-query";
import type { ApiError } from "../api/client";
import {
  getFetchCurrentUserQueryKey,
  useFetchCurrentUser,
  useLoginUser,
  useLogoutUser,
  useResetUserPassword,
} from "../api/endpoints";

const isAccountExpiredError = (error: unknown) => {
  const { status, info } = (error ?? {}) as ApiError;
  return status === 403 && info?.detail === "account expired";
};

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [resetTime, setResetTime] = useState<string | null>(null);
  const [expiresAt, setExpiresAt] = useState<string | null>(null);
  const [userID, setUserID] = useState<string | null>(null);
  const [counselorID, setCounselorID] = useState<string | null>(null);
  const [studentID, setStudentID] = useState<string | null>(null);
  const [error, setError] = useState<ApiError | null>(null);
  const [path, setPath] = useState("/");
  const [serverExpired, setServerExpired] = useState(false);
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const queryClient = useQueryClient();

  const { mutateAsync: loginFunc } = useLoginUser();
  const { mutateAsync: logoutFunc } = useLogoutUser();
  const { mutateAsync: passFunc } = useResetUserPassword();

  const clearSession = () => {
    setResetTime(null);
    setExpiresAt(null);
    setUserID(null);
    setCounselorID(null);
    setStudentID(null);
    setServerExpired(false);
    queryClient.removeQueries({ queryKey: getFetchCurrentUserQueryKey() });
  };

  // restore the session on page load from the session cookie
  const { data: currentUser } = useFetchCurrentUser({
    query: { retry: false, refetchOnWindowFocus: false },
  });

  useEffect(() => {
    if (!currentUser) {
      return;
    }
    setResetTime(currentUser.reset_time);
    setExpiresAt(currentUser.expires_at);
    setUserID(currentUser.id);
    setCounselorID(currentUser.counselor_id);
    setStudentID(currentUser.student_id);
  }, [currentUser]);

  const login = (email: string, password: string) => {
    return loginFunc(
      { data: { email, password } },
      {
        onSuccess: (user) => {
          setResetTime(user.reset_time);
          setExpiresAt(user.expires_at);
          setUserID(user.id);
          setCounselorID(user.counselor_id);
          setStudentID(user.student_id);
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
    return logoutFunc(undefined, {
      onSettled: clearSession,
      onSuccess: () => {
        void navigate(path);
        setPath("/");
      },
      onError: (err) => {
        setError(err);
      },
    });
  };

  useEffect(() => {
    const onEvent = (
      event: QueryCacheNotifyEvent | MutationCacheNotifyEvent,
    ) => {
      if (
        event.type === "updated" &&
        event.action.type === "error" &&
        isAccountExpiredError(event.action.error)
      ) {
        setServerExpired(true);
      }
    };
    const unsubscribeQueries = queryClient.getQueryCache().subscribe(onEvent);
    const unsubscribeMutations = queryClient
      .getMutationCache()
      .subscribe(onEvent);
    return () => {
      unsubscribeQueries();
      unsubscribeMutations();
    };
  }, [queryClient]);

  const new_password = (new_password: string) => {
    return passFunc(
      { data: { new_password } },
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

  const isExpired =
    serverExpired ||
    (expiresAt !== null && new Date(expiresAt).getTime() <= Date.now());

  return (
    <AuthContext
      value={{
        resetTime: resetTime,
        expiresAt: expiresAt,
        userID: userID,
        counselorID: counselorID,
        studentID: studentID,
        login: login,
        logout: logout,
        reset_password: new_password,
      }}
    >
      {isExpired && pathname !== ACCOUNT_EXPIRED_PATH ? (
        <Navigate to={ACCOUNT_EXPIRED_PATH} replace />
      ) : !isExpired && pathname === ACCOUNT_EXPIRED_PATH ? (
        <Navigate to="/" replace />
      ) : (
        children
      )}
    </AuthContext>
  );
}
