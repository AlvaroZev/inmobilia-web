import type { CostResult } from './cost';

export interface BillOfMaterials {
  furnitureId: string;
  furnitureName?: string;
  generatedAt: string;
  parts: BOMPartLine[];
  hardware: BOMHardwareLine[];
  edgeBanding: BOMEdgeLine[];
  cost?: CostResult;
  summary: BOMSummary;
}

export interface BOMPartLine {
  partId: string;
  name: string;
  type: string;
  volumeId: string;
  width: number;
  height: number;
  thickness: number;
  materialId: string;
  materialName: string;
  grainDirection?: string;
  areaM2: number;
}

export interface BOMHardwareLine {
  hardwareType: string;
  name: string;
  quantity: number;
  unitCost?: number;
  total?: number;
}

export interface BOMEdgeLine {
  material: string;
  totalLengthM: number;
}

export interface BOMSummary {
  partCount: number;
  hardwareCount: number;
  totalBoardM2: number;
  totalEdgeM: number;
}

export interface CutPlan {
  furnitureId: string;
  furnitureName?: string;
  generatedAt: string;
  sheets: CutSheet[];
}

export interface CutSheet {
  materialId: string;
  materialName: string;
  thickness: number;
  parts: CutPartLine[];
  totalAreaM2: number;
}

export interface CutPartLine {
  partId: string;
  name: string;
  width: number;
  height: number;
  thickness: number;
  grain?: string;
  quantity: number;
}
