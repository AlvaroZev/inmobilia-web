import { create } from 'zustand';
import type { ResolvedFurniture } from '@/domain/resolved-furniture';
import { normalizeResolvedFurniture, parseResolvedFurnitureJson } from '@/utils/load-resolved';

interface ViewerState {
  resolvedFurniture: ResolvedFurniture | null;
  overrideJsonText: string;
  overrideFurniture: ResolvedFurniture | null;
  overrideError: string | null;
  selectedVolumeId: string | null;
  openDrawers: Record<string, 'closed' | 'half'>;
  showVolumes: boolean;
  showFeatures: boolean;
  showFronts: boolean;
  setResolvedFurniture: (furniture: ResolvedFurniture) => void;
  clearResolvedFurniture: () => void;
  setOverrideJsonText: (text: string) => void;
  applyJsonOverride: () => boolean;
  clearJsonOverride: () => void;
  setSelectedVolumeId: (id: string | null) => void;
  toggleDrawer: (drawerId: string) => void;
  clearOpenDrawers: () => void;
  toggleVolumes: () => void;
  toggleFeatures: () => void;
  toggleFronts: () => void;
}

export const useViewerStore = create<ViewerState>((set, get) => ({
  resolvedFurniture: null,
  overrideJsonText: '',
  overrideFurniture: null,
  overrideError: null,
  selectedVolumeId: null,
  openDrawers: {},
  showVolumes: false,
  showFeatures: true,
  showFronts: true,
  setResolvedFurniture: (furniture) =>
    set({
      resolvedFurniture: normalizeResolvedFurniture(furniture),
      selectedVolumeId: null,
      openDrawers: {},
    }),
  clearResolvedFurniture: () => set({ resolvedFurniture: null, selectedVolumeId: null, openDrawers: {} }),
  setOverrideJsonText: (text) => set({ overrideJsonText: text, overrideError: null }),
  applyJsonOverride: () => {
    const text = get().overrideJsonText;
    if (!text.trim()) {
      set({ overrideFurniture: null, overrideError: null, selectedVolumeId: null, openDrawers: {} });
      return true;
    }
    const result = parseResolvedFurnitureJson(text);
    if (!result.ok) {
      set({ overrideError: result.error });
      return false;
    }
    set({
      overrideFurniture: result.furniture,
      overrideError: null,
      selectedVolumeId: null,
      openDrawers: {},
    });
    return true;
  },
  clearJsonOverride: () =>
    set({
      overrideJsonText: '',
      overrideFurniture: null,
      overrideError: null,
      selectedVolumeId: null,
      openDrawers: {},
    }),
  setSelectedVolumeId: (id) => set({ selectedVolumeId: id }),
  toggleDrawer: (drawerId) =>
    set((state) => ({
      openDrawers: {
        ...state.openDrawers,
        [drawerId]: state.openDrawers[drawerId] === 'half' ? 'closed' : 'half',
      },
    })),
  clearOpenDrawers: () => set({ openDrawers: {} }),
  toggleVolumes: () => set((state) => ({ showVolumes: !state.showVolumes })),
  toggleFeatures: () => set((state) => ({ showFeatures: !state.showFeatures })),
  toggleFronts: () => set((state) => ({ showFronts: !state.showFronts })),
}));
