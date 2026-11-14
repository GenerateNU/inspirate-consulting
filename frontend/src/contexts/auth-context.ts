import { createContext } from "react";

import type { CurrentUserRole, SessionUser } from "../api/models";

export const ACCOUNT_EXPIRED_PATH = "/account-expired";
export const CLAIM_ACCOUNT_PATH = "/claim-account";
export const LOGIN_PATH = "/login";
// Pages reachable without being logged in.
export const PUBLIC_PATHS = [LOGIN_PATH, "/student/login", "/counselor/login"];

export interface AuthContextType {
  resetTime: string | null;
  expiresAt: string | null;
  counselorID: string | null;
  studentID: string | null;
  userID: string | null;
  role: CurrentUserRole | null;
  login: (email: string, password: string) => Promise<SessionUser>;
  logout: () => Promise<void>;
  reset_password: (new_password: string) => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
