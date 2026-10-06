import { createContext } from "react";
import type {
    CreateUserBody,
    CreateUserInputBodyRole,
    
} from '../api/models'
import type { createUserResponse } from "../api/endpoints";

export interface AuthContextType {
    resetTime: string | null,
    login: (
        email: string,
        password: string
    ) => void;
    create: (
        name: string,
        email: string,
        pfp_key: string | null,
        role: CreateUserInputBodyRole,
    ) => Promise<createUserResponse>;
    force_password_reset: (
        supabase_id: string
    ) => void;
    logout: () => void;
    reset_password: (
        new_password: string,
    ) => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)