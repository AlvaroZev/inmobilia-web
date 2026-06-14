import type { ManufacturingModel, Part } from '@/domain/manufacturing';
import { EPSILON } from './geometry-math';
import type { ValidationError, ValidationResult } from './validation';

function pushError(errors: ValidationError[], field: string, message: string): void {
  errors.push({ field, message });
}

export function countPartsByType(parts: Part[], partType: string): number {
  return parts.filter((part) => part.type === partType).length;
}

export function totalBoardAreaM2(parts: Part[]): number {
  return parts.reduce((sum, part) => {
    const areaMm2 = part.width * part.height;
    return sum + areaMm2 / 1_000_000;
  }, 0);
}

export function totalEdgeBandingM(model: ManufacturingModel): number {
  return model.edgeBanding.reduce((sum, edge) => sum + edge.length, 0) / 1000;
}

export function validateManufacturingModel(model: ManufacturingModel): ValidationResult {
  const errors: ValidationError[] = [];

  if (!model.furnitureId) pushError(errors, 'furnitureId', 'furnitureId is required');
  if (model.parts.length === 0) pushError(errors, 'parts', 'at least one part is required');

  const partIds = new Map<string, string>();
  model.parts.forEach((part, i) => {
    const prefix = `parts[${i}]`;
    if (!part.id) pushError(errors, `${prefix}.id`, 'id is required');
    if (part.id) {
      const existing = partIds.get(part.id);
      if (existing) {
        pushError(errors, `${prefix}.id`, `duplicate id "${part.id}" (also used in ${existing})`);
      } else {
        partIds.set(part.id, prefix);
      }
    }
    if (part.width <= EPSILON || part.height <= EPSILON || part.thickness <= EPSILON) {
      pushError(errors, prefix, 'part dimensions must be greater than zero');
    }
    if (!part.volumeId) pushError(errors, `${prefix}.volumeId`, 'volumeId is required');
  });

  model.edgeBanding.forEach((edge, i) => {
    if (!edge.partId) pushError(errors, `edgeBanding[${i}].partId`, 'partId is required');
    if (edge.length <= EPSILON) {
      pushError(errors, `edgeBanding[${i}].length`, 'length must be greater than zero');
    }
  });

  model.hardware.forEach((hw, i) => {
    if (!hw.type) pushError(errors, `hardware[${i}].type`, 'type is required');
    if (hw.quantity <= 0) pushError(errors, `hardware[${i}].quantity`, 'quantity must be greater than zero');
  });

  return { valid: errors.length === 0, errors };
}
