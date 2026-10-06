import { createContext } from "react";
import type {
    CreateUserBody,
} from '../api/models'

export interface AuthContextType {
    resetTime: string | null,
    login: (
        email: string,
        password: string
    ) => void;
    create: (
        name: string,
        email: string,
        role: string
    ) => CreateUserBody;
    force_password_reset: (
        supabase_id: string
    ) => void;
    logout: () => void;
    reset_password: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)