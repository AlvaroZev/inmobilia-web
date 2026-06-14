/**
 * Cost Engine output.
 */

export interface CostResult {
  furnitureId: string;
  currency: string;
  materials: MaterialCostLine[];
  hardware: HardwareCostLine[];
  labor: LaborCost;
  waste: WasteCost;
  subtotal: number;
  total: number;
}

export interface MaterialCostLine {
  materialId: string;
  name: string;
  areaM2: number;
  unitCostPerM2: number;
  total: number;
}

export interface HardwareCostLine {
  hardwareType: string;
  name: string;
  quantity: number;
  unitCost: number;
  total: number;
}

export interface LaborCost {
  hours: number;
  ratePerHour: number;
  total: number;
}

export interface WasteCost {
  areaM2: number;
  percentage: number;
  total: number;
}
