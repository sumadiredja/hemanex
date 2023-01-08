// Utilities
import { defineStore } from "pinia";

export const useGlobalStore = defineStore("app", {
  state: () => ({
    NavbarCollapse: true,
  }),
  actions: {
    changeCollapse() {
      this.NavbarCollapse = !this.NavbarCollapse;
    },
    isCollapse() {
      if (this.NavbarCollapse) {
        return true;
      }
      return false;
    },
  },
});
