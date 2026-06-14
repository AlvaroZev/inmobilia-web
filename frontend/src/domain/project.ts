import type { CostResult } from './cost';
import type { FurnitureDefinition } from './furniture-definition';
import type { ManufacturingModel } from './manufacturing';
import type { ResolvedFurniture } from './resolved-furniture';
import type { RoomSetupValues } from '@/utils/room-setup';

export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
}

export interface PipelineSnapshot {
  furniture: FurnitureDefinition;
  resolved?: ResolvedFurniture;
  manufacturing?: ManufacturingModel;
  cost?: CostResult;
}

export interface SavedProject {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
  messages: ChatMessage[];
  roomSetup: RoomSetupValues;
  runPipeline: boolean;
  result: PipelineSnapshot;
}
