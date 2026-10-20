import { createContext } from "react";

import type {
  loginUserResponse,
  logoutUserResponse,
  resetUserPasswordResponse,
} from "../api/endpoints";

export interface AuthContextType {
  resetTime: string | null;
  counselorID: string | null;
  studentID: string | null;
  userID: string | null;
  login: (email: string, password: string) => Promise<loginUserResponse>;
  logout: () => Promise<logoutUserResponse>;
  reset_password: (new_password: string) => Promise<resetUserPasswordResponse>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
