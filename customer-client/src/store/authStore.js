import { create } from "zustand"

export const useAuthStore = create((set) => ({
    // isAuth: false,
    // setIsAuth: (isAuth) => set({ isAuth: isAuth })
    shopId: "",
    setShopId: (shopId) => set({ shopId: shopId }),
    tableId: "",
    setTableId: (tableId) => set({ tableId: tableId }),
    placeOrderItems: [],
    setPlaceOrderItems: (items) => set({ placeOrderItems: items })
}))