import type { ResolvedFurniture, ResolvedVolume } from '@/domain/resolved-furniture';
import { validateResolvedFurniture } from '@/utils/resolved-furniture';

function normalizeVolume(raw: ResolvedVolume): ResolvedVolume {
  return {
    ...raw,
    children: (raw.children ?? []).map(normalizeVolume),
    features: raw.features ?? [],
    fronts: raw.fronts ?? [],
  };
}

export function normalizeResolvedFurniture(raw: ResolvedFurniture): ResolvedFurniture {
  return {
    ...raw,
    root: normalizeVolume(raw.root),
  };
}

export type ParseResolvedResult =
  | { ok: true; furniture: ResolvedFurniture }
  | { ok: false; error: string };

export function parseResolvedFurnitureJson(text: string): ParseResolvedResult {
  const trimmed = text.trim();
  if (!trimmed) {
    return { ok: false, error: 'El JSON está vacío' };
  }

  let raw: ResolvedFurniture;
  try {
    raw = JSON.parse(trimmed) as ResolvedFurniture;
  } catch (err) {
    const message = err instanceof Error ? err.message : 'JSON inválido';
    return { ok: false, error: message };
  }

  const furniture = normalizeResolvedFurniture(raw);
  const validation = validateResolvedFurniture(furniture);
  if (!validation.valid) {
    const detail = validation.errors.map((e) => `${e.field}: ${e.message}`).join('; ');
    return { ok: false, error: detail };
  }

  return { ok: true, furniture };
}
