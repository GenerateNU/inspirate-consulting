import { createContext } from "react";

// login should redirect conditionally based on the user field
export interface AuthContextType {
    isAuthenticated: boolean
    resetTime: Date,
    login: (
        email: string,
        password: string
    ) => void;
    create: (
        name: string,
        email: string,
        username: string,
        role: string
    ) => void;
    force_password_reset: (
        supabase_id: string
    ) => void;
    logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)