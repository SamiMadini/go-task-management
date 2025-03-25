import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux"
import type { RootState, AppDispatch } from "./store"

export const useAppDispatch: () => AppDispatch = useDispatch
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector

export const useAuth = () => {
  const { user, accessToken, refreshToken, isLoading, error } = useAppSelector((state) => state.auth)

  return {
    user,
    accessToken,
    refreshToken,
    isLoading,
    error,
    isAuthenticated: !!user && !!accessToken,
  }
}
