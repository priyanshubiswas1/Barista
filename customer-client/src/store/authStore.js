import { create } from "zustand"

export const useAuthStore = create((set) => ({
    // isAuth: false,
    // setIsAuth: (isAuth) => set({ isAuth: isAuth })
    placeOrderItems: [],
    setPlaceOrderItems: (items) => set({ placeOrderItems: items })
}))