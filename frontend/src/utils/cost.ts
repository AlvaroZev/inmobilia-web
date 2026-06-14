import type { CostResult } from '@/domain/cost';
import type { ValidationError, ValidationResult } from './validation';

function pushError(errors: ValidationError[], field: string, message: string): void {
  errors.push({ field, message });
}

export function sumMaterialTotal(result: CostResult): number {
  return result.materials.reduce((sum, line) => sum + line.total, 0);
}

export function sumHardwareTotal(result: CostResult): number {
  return result.hardware.reduce((sum, line) => sum + line.total, 0);
}

export function validateCostResult(result: CostResult): ValidationResult {
  const errors: ValidationError[] = [];

  if (!result.furnitureId) pushError(errors, 'furnitureId', 'furnitureId is required');
  if (!result.currency) pushError(errors, 'currency', 'currency is required');
  if (result.total < 0) pushError(errors, 'total', 'total must be non-negative');

  const expectedSubtotal =
    sumMaterialTotal(result) + sumHardwareTotal(result) + result.labor.total;

  if (Math.abs(result.subtotal - expectedSubtotal) > 0.02) {
    pushError(
      errors,
      'subtotal',
      `subtotal ${result.subtotal} does not match line items ${expectedSubtotal}`,
    );
  }

  const expectedTotal = result.subtotal + result.waste.total;
  if (Math.abs(result.total - expectedTotal) > 0.02) {
    pushError(
      errors,
      'total',
      `total ${result.total} does not match subtotal + waste ${expectedTotal}`,
    );
  }

  return { valid: errors.length === 0, errors };
}
