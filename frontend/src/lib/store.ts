import { configureStore } from "@reduxjs/toolkit"
import { persistStore, FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER } from "redux-persist"
import storage from "redux-persist/lib/storage"
import authReducer from "./features/auth/authSlice"
import { persistReducer } from "redux-persist"
import { combineReducers } from "@reduxjs/toolkit"
import { AuthState } from "./features/auth/authSlice"

type RootState = {
  auth: AuthState
}

const persistConfig = {
  key: "root",
  storage,
  whitelist: ["auth"], // Only persist auth state
  blacklist: [], // Don't persist these reducers
  debug: true, // Always enable debug logging
  stateReconciler: (inboundState: Partial<RootState>, originalState: RootState, reducedState: RootState) => {
    console.log("State reconciler called:", {
      inboundState: {
        hasUser: !!inboundState.auth?.user,
        hasAccessToken: !!inboundState.auth?.accessToken,
        hasRefreshToken: !!inboundState.auth?.refreshToken,
      },
      originalState: {
        hasUser: !!originalState.auth.user,
        hasAccessToken: !!originalState.auth.accessToken,
        hasRefreshToken: !!originalState.auth.refreshToken,
      },
      reducedState: {
        hasUser: !!reducedState.auth.user,
        hasAccessToken: !!reducedState.auth.accessToken,
        hasRefreshToken: !!reducedState.auth.refreshToken,
      },
    })

    // Ensure we're properly merging the persisted state
    return {
      ...reducedState,
      auth: {
        ...reducedState.auth,
        ...inboundState.auth,
      },
    }
  },
}

const rootReducer = combineReducers({
  auth: authReducer,
})

const persistedReducer = persistReducer(persistConfig, rootReducer)

export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
      },
    }),
  devTools: true, // Always enable dev tools
})

export const persistor = persistStore(store)

// Add debug logging for store state changes
store.subscribe(() => {
  const state = store.getState() as RootState
  console.log("Redux store state changed:", {
    hasUser: !!state.auth.user,
    hasAccessToken: !!state.auth.accessToken,
    hasRefreshToken: !!state.auth.refreshToken,
    accessToken: state.auth.accessToken?.substring(0, 10) + "...",
    refreshToken: state.auth.refreshToken?.substring(0, 10) + "...",
    user: state.auth.user,
  })
})

export type { RootState }
export type AppDispatch = typeof store.dispatch
