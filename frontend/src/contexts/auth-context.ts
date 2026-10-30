import { createContext } from "react";

import type { SessionUser } from "../api/models";

export const ACCOUNT_EXPIRED_PATH = "/account-expired";

export interface AuthContextType {
  resetTime: string | null;
  expiresAt: string | null;
  counselorID: string | null;
  studentID: string | null;
  userID: string | null;
  login: (email: string, password: string) => Promise<SessionUser>;
  logout: () => Promise<void>;
  reset_password: (new_password: string) => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
