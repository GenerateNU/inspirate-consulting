import type {
    CreateUserBody,
    LoginResponse
} from '../api/models'
import { AuthContext } from './auth-context'
import { useState} from "react";
import { useNavigate } from "react-router-dom";
import { useLoginUser, useLogoutUser, useResetUserPassword, useCreateUser, useForcePasswordReset, forcePasswordReset} from '../api/endpoints';

export function AuthProvider({children}: {children: React.ReactNode}) {
    const [resetTime, setResetTime] = useState<string | null >('')
    const [supabaseId, setSupabaseID] = useState('')
    const navigate = useNavigate();

    const {trigger: loginFunc} = useLoginUser();
    const createFunc = useCreateUser();
    const logoutFunc = useLogoutUser();
    const passFunc = useResetUserPassword();
    const resetFunc  = useForcePasswordReset(supabaseId);

    const login = (
        email: string,
        password: string,
    ) => {
        loginFunc(
            {email, password},
            {
                onSuccess: (resp: LoginResponse) => {
                    navigate("/")
                }
            }
        );
    }

    const logout = () => {
        logoutFunc(
            {
                onSuccess: () => {
                    navigate("/login")
                }
            }
        );
    }

    const create = (
        name: string,
        email: string,
        role: string,
    ) => {
        createFunc(
            {data: {name, email, role}}
        );

    }

    const force_password_reset = (
        supabase_id: string,
    ) => {
        // set
        

        setResetTime("one second ago")
    }

    const reset_password = () => {
        passFunc(
             {
                onSuccess: () => {
                    navigate("/login")
                }
            }
        );
    }


    return (
        <AuthContext value={{
            resetTime: resetTime,
            login: login,
            create: create,
            force_password_reset: force_password_reset,
            logout: logout,
            reset_password: reset_password
        }}>
        {children}
        </AuthContext>
    )
}