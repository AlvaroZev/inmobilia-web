import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { SavedProject } from '@/domain/project';

const STORAGE_KEY = 'inmobilia-projects';
const MAX_PROJECTS = 50;

export interface UpsertProjectInput {
  id?: string;
  name: string;
  description: string;
  messages: SavedProject['messages'];
  roomSetup: SavedProject['roomSetup'];
  runPipeline: boolean;
  result: SavedProject['result'];
}

interface ProjectsState {
  projects: SavedProject[];
  upsertProject: (input: UpsertProjectInput) => string;
  deleteProject: (id: string) => void;
}

export const useProjectsStore = create<ProjectsState>()(
  persist(
    (set, get) => ({
      projects: [],
      upsertProject: (input) => {
        const now = new Date().toISOString();
        const id = input.id ?? crypto.randomUUID();
        const existing = get().projects.find((project) => project.id === id);
        const project: SavedProject = {
          id,
          name: input.name,
          description: input.description,
          createdAt: existing?.createdAt ?? now,
          updatedAt: now,
          messages: input.messages,
          roomSetup: input.roomSetup,
          runPipeline: input.runPipeline,
          result: input.result,
        };

        set((state) => ({
          projects: [project, ...state.projects.filter((entry) => entry.id !== id)].slice(0, MAX_PROJECTS),
        }));

        return id;
      },
      deleteProject: (id) =>
        set((state) => ({
          projects: state.projects.filter((project) => project.id !== id),
        })),
    }),
    { name: STORAGE_KEY },
  ),
);
