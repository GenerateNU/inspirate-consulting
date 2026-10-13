import { createContext } from "react";

import type { User } from "../api/models";
import { LoginUserResponse } from "../api/schemas/inspirateConsultingAPI";
import { useLoginUser, useLogoutUser, useResetUserPassword, type loginUserResponse} from '../api/endpoints';

export const ACCOUNT_EXPIRED_PATH = "/account-expired";

export interface AuthContextType {
  resetTime: string | null;
  expiresAt: string | null;
  userID: string | null;
  login: (
        email: string,
        password: string
    ) => Promise<LoginUserResponse>;
    logout: () => Promise<logoutUserResponse>;
    reset_password: (
        new_password: string,
    ) => Promise<resetUserPasswordResponse>;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
