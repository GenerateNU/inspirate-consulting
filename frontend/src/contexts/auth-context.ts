import { createContext } from "react";

import type { User } from "../api/models";

export interface AuthContextType {
  resetTime: string | null;
  userID: string | null;
  counselorID: string | null;
  studentID: string | null;
  login: (email: string, password: string) => Promise<loginUserResponse>;
  logout: () => Promise<logoutUserResponse>;
  reset_password: (new_password: string) => Promise<resetUserPasswordResponse>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
