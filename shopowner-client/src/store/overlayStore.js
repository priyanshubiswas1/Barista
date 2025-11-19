import { create } from "zustand"

export const useOverlay = create((set) => ({
    overlay: "",
    overlayData: null,
    setOverlay: (overlay, data = null) => set({ overlay: overlay, overlayData: data })
}))