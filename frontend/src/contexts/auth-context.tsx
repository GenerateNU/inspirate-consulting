import type {
    CreateUserBody,
    CreateUserInputBodyRole,
    LoginResponse
} from '../api/models'
import { AuthContext } from './auth-context'
import { useState} from "react";
import { useNavigate } from "react-router-dom";
import { useLoginUser, useLogoutUser, useResetUserPassword, useCreateUser, useForcePasswordReset, forcePasswordReset, loginUserResponse, resetUserPasswordResponse, resetUserPassword, createUserResponse, createUserResponseError} from '../api/endpoints';

export function AuthProvider({children}: {children: React.ReactNode}) {
    const [resetTime, setResetTime] = useState<string | null >(null)
    const [supabaseId, setSupabaseID] = useState('')
    const navigate = useNavigate();

    const {trigger: loginFunc} = useLoginUser();
    const {trigger: createFunc} = useCreateUser();
    const {trigger: logoutFunc} = useLogoutUser();
    const {trigger: passFunc} = useResetUserPassword();
    const {trigger: resetFunc}  = useForcePasswordReset(supabaseId);

    const login = (
        email: string,
        password: string,
    ) => {
        loginFunc(
            {email, password},
            {
                onSuccess: ()  => {
                    navigate("/")
                },

                onError: () => {
                    console.log("did not work")
                }
            }
        );
    }

    const logout = () => {
    }
    logoutFunc(
        {},
        {
            onSuccess: () => {
                navigate("/")
            }
        }
    );

    const create = (
        name: string,
        email: string,
        pfp_key: string | null,
        role: CreateUserInputBodyRole,
    ) => {
        return createFunc(
            {name, email, pfp_key, role},
            {
                onSuccess: (resp: createUserResponse) => {
                    return resp
                },

            }
        );
    }

    const new_password = (
        new_password: string
    ) => {
        passFunc(
            {new_password},
            {
                onSuccess: () => {
                     navigate("/")
                },

                onError: () => {
                    console.log("resetting password did not work")
                }
            }
        );
    }

    const force_password_reset = () => {
    }
    resetFunc(
        {supabaseId},
        {
            onSuccess: () => {
                navigate("/")
            }
        }
    );

    return (
        <AuthContext value={{
            resetTime: resetTime,
            login: login,
            create: create,
            force_password_reset: force_password_reset,
            logout: logout,
            reset_password: new_password
        }}>
        {children}
        </AuthContext>
    )
}