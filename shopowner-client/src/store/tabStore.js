import { create } from "zustand"
import { Tabs } from "../constants/constants"

export const useTabs = create((set) => ({
    tab: Tabs.MENU,
    setTab: (tab) => set({ tab: tab })
}))