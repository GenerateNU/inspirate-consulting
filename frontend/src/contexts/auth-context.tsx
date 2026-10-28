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
  getGetCurrentUserQueryKey,
  useGetCurrentUser,
  useLoginUser,
  useLogoutUser,
  useResetUserPassword,
} from "../api/endpoints";

const isAccountExpiredError = (error: unknown) => {
  const { status, info } = (error ?? {}) as ApiError;
  return status === 403 && info?.detail === "account expired";
};

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [error, setError] = useState<ApiError | null>(null);
  const [path, setPath] = useState("/");
  const [serverExpired, setServerExpired] = useState(false);
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const queryClient = useQueryClient();

  const { data: currentUser } = useGetCurrentUser();
  const resetTime = currentUser?.reset_time ?? null;
  const expiresAt = currentUser?.expires_at ?? null;
  const userID = currentUser?.id ?? null;
  const counselorID = currentUser?.counselor_id ?? null;
  const studentID = currentUser?.student_id ?? null;

  const { mutateAsync: loginFunc } = useLoginUser();
  const { mutateAsync: logoutFunc } = useLogoutUser();
  const { mutateAsync: passFunc } = useResetUserPassword();

  const clearSession = () => {
    void queryClient.resetQueries({ queryKey: getGetCurrentUserQueryKey() });
    setServerExpired(false);
  };

  const login = (email: string, password: string) => {
    return loginFunc(
      { data: { email, password } },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: getGetCurrentUserQueryKey(),
          });
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
