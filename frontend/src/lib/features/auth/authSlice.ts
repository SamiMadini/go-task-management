import { createSlice, PayloadAction } from "@reduxjs/toolkit"

export interface User {
  id: string
  email: string
  handle: string
}

export interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isLoading: boolean
  error: string | null
}

const initialState: AuthState = {
  user: null,
  accessToken: null,
  refreshToken: null,
  isLoading: false,
  error: null,
}

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (
      state,
      action: PayloadAction<{
        user: User
        accessToken: string
        refreshToken: string
      }>
    ) => {
      console.log("Setting credentials in auth slice:", {
        hasUser: !!action.payload.user,
        hasAccessToken: !!action.payload.accessToken,
        hasRefreshToken: !!action.payload.refreshToken,
        user: action.payload.user,
        accessToken: action.payload.accessToken.substring(0, 10) + "...",
        refreshToken: action.payload.refreshToken.substring(0, 10) + "...",
      })
      state.user = { ...action.payload.user }
      state.accessToken = action.payload.accessToken
      state.refreshToken = action.payload.refreshToken
      state.error = null
    },
    setAccessToken: (state, action: PayloadAction<string>) => {
      console.log("Setting access token in auth slice:", {
        hasToken: !!action.payload,
        tokenPreview: action.payload.substring(0, 10) + "...",
        currentUser: state.user,
      })
      state.accessToken = action.payload
    },
    setError: (state, action: PayloadAction<string>) => {
      console.log("Setting error in auth slice:", {
        error: action.payload,
        currentUser: state.user,
        hasAccessToken: !!state.accessToken,
      })
      state.error = action.payload
      state.isLoading = false
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      console.log("Setting loading state in auth slice:", {
        isLoading: action.payload,
        currentUser: state.user,
        hasAccessToken: !!state.accessToken,
      })
      state.isLoading = action.payload
    },
    logout: (state) => {
      console.log("Logging out user from auth slice:", {
        previousUser: state.user,
        hadAccessToken: !!state.accessToken,
        hadRefreshToken: !!state.refreshToken,
      })
      state.user = null
      state.accessToken = null
      state.refreshToken = null
      state.error = null
    },
  },
})

export const { setCredentials, setAccessToken, setError, setLoading, logout } = authSlice.actions
export default authSlice.reducer
