import { reactive } from 'vue';
import type { Kid, Behavior, Prize } from '../types';
import * as api from '../api';

// Simple reactive store using Vue's reactive directly
const state = reactive({
  kids: [] as Kid[],
  behaviors: [] as Behavior[],
  prizes: [] as Prize[],
  currentKid: null as Kid | null,
  loading: false,
  error: null as string | null,

  async fetchKids() {
    state.loading = true;
    state.error = null;
    try {
      const kids = await api.getKids();
      state.kids = kids;
      // Update currentKid if it exists in the new kids list
      if (state.currentKid) {
        const updatedKid = kids.find(k => k.id === state.currentKid?.id);
        if (updatedKid) {
          state.currentKid = { ...updatedKid };
        } else {
          state.currentKid = kids[0] || null;
        }
      } else {
        state.currentKid = kids[0] || null;
      }
      state.loading = false;
    } catch (e) {
      state.error = (e as Error).message;
      state.loading = false;
    }
  },

  async fetchBehaviors() {
    try {
      state.behaviors = await api.getBehaviors();
    } catch (e) {
      state.error = (e as Error).message;
    }
  },

  setCurrentKid(kid: Kid | null) {
    state.currentKid = kid;
  },

  async addPoints(behaviorId: number) {
    if (!state.currentKid) return;
    try {
      const updatedKid = await api.addPoints(state.currentKid.id, behaviorId);
      // Update kids array with updated kid
      const idx = state.kids.findIndex(k => k.id === updatedKid.id);
      if (idx !== -1) {
        state.kids[idx] = updatedKid;
      }
      // Update currentKid reference to trigger reactivity
      state.currentKid = { ...updatedKid };
    } catch (e) {
      state.error = (e as Error).message;
    }
  },

  async subPoints(behaviorId: number) {
    if (!state.currentKid) return;
    try {
      const updatedKid = await api.subPoints(state.currentKid.id, behaviorId);
      const idx = state.kids.findIndex(k => k.id === updatedKid.id);
      if (idx !== -1) {
        state.kids[idx] = updatedKid;
      }
      // Create new object reference to trigger Vue reactivity
      state.currentKid = { ...updatedKid };
    } catch (e) {
      state.error = (e as Error).message;
    }
  },
});

export function useStore() {
  return state;
}

(window as any).__store = state;
