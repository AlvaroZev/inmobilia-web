import { create } from 'zustand';
import type { ChatMessage, PipelineSnapshot } from '@/domain/project';
import { defaultRoomSetupValues, type RoomSetupValues } from '@/utils/room-setup';

export const welcomeMessage: ChatMessage = {
  id: 'welcome',
  role: 'assistant',
  content:
    'Describe el mueble en lenguaje natural. Generaré un Volume Tree (FurnitureDefinition) y opcionalmente el pipeline completo.',
};

interface AssistantState {
  messages: ChatMessage[];
  input: string;
  name: string;
  runPipeline: boolean;
  result: PipelineSnapshot | null;
  error: string | null;
  roomSetup: RoomSetupValues;
  roomSetupErrors: string[];
  activeProjectId: string | null;
  saveNotice: string | null;
  setInput: (input: string) => void;
  setName: (name: string) => void;
  setRunPipeline: (runPipeline: boolean) => void;
  setResult: (result: PipelineSnapshot | null) => void;
  setError: (error: string | null) => void;
  setRoomSetup: (roomSetup: RoomSetupValues) => void;
  setRoomSetupErrors: (errors: string[]) => void;
  setActiveProjectId: (id: string | null) => void;
  setSaveNotice: (notice: string | null) => void;
  setMessages: (messages: ChatMessage[] | ((prev: ChatMessage[]) => ChatMessage[])) => void;
  appendMessage: (message: ChatMessage) => void;
  loadSession: (session: {
    messages: ChatMessage[];
    name: string;
    roomSetup: RoomSetupValues;
    runPipeline: boolean;
    result: PipelineSnapshot;
    activeProjectId: string;
  }) => void;
  resetSession: () => void;
}

export const useAssistantStore = create<AssistantState>((set) => ({
  messages: [welcomeMessage],
  input: '',
  name: '',
  runPipeline: true,
  result: null,
  error: null,
  roomSetup: defaultRoomSetupValues(),
  roomSetupErrors: [],
  activeProjectId: null,
  saveNotice: null,
  setInput: (input) => set({ input }),
  setName: (name) => set({ name }),
  setRunPipeline: (runPipeline) => set({ runPipeline }),
  setResult: (result) => set({ result }),
  setError: (error) => set({ error }),
  setRoomSetup: (roomSetup) => set({ roomSetup }),
  setRoomSetupErrors: (roomSetupErrors) => set({ roomSetupErrors }),
  setActiveProjectId: (activeProjectId) => set({ activeProjectId }),
  setSaveNotice: (saveNotice) => set({ saveNotice }),
  setMessages: (messages) =>
    set((state) => ({
      messages: typeof messages === 'function' ? messages(state.messages) : messages,
    })),
  appendMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  loadSession: (session) =>
    set({
      activeProjectId: session.activeProjectId,
      messages: session.messages,
      name: session.name,
      roomSetup: session.roomSetup,
      runPipeline: session.runPipeline,
      result: session.result,
      error: null,
      saveNotice: null,
      roomSetupErrors: [],
      input: '',
    }),
  resetSession: () =>
    set({
      activeProjectId: null,
      messages: [welcomeMessage],
      name: '',
      input: '',
      result: null,
      error: null,
      saveNotice: null,
      roomSetup: defaultRoomSetupValues(),
      roomSetupErrors: [],
    }),
}));
