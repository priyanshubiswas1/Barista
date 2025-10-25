import { create } from "zustand"

export const useAuthStore = create((set) => ({
    // isAuth: false,
    // setIsAuth: (isAuth) => set({ isAuth: isAuth })
    orderItems: [],
    setOrderItems: (items) => set({ orderItems: items })
}))