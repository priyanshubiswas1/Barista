import { create } from "zustand"

export const useOverlay = create((set) => ({
    overlay: "",
    setOverlay: (overlay) => set({ overlay: overlay })
}))