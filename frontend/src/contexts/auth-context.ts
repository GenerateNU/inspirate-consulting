import { createContext } from "react";

import type { User } from "../api/models";

export interface AuthContextType {
  resetTime: string | null;
  userID: string | null;
  login: (email: string, password: string) => Promise<User>;
  logout: () => Promise<void>;
  reset_password: (new_password: string) => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
