/** Redondeo a milímetros enteros: ≤0,4 baja; ≥0,5 sube. */
export function roundMm(value: number): number {
  return Math.round(value);
}
