type RunnerSide = 'left' | 'right';

export interface RunnerPanelSpec {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  role: 'runner_cabinet' | 'runner_drawer';
}

/** Perfil en U: 3 paneles (web + labio inferior + labio superior). */
function uChannelPanels(
  zoneStartX: number,
  zoneWidth: number,
  runnerY: number,
  runnerHeight: number,
  runnerDepth: number,
  runnerZ: number,
  role: 'runner_cabinet' | 'runner_drawer',
  side: RunnerSide,
  lipScale: number,
): RunnerPanelSpec[] {
  if (zoneWidth <= 0 || runnerHeight <= 0 || runnerDepth <= 0) {
    return [];
  }

  const lipH = Math.min(1.4, runnerHeight * 0.12);
  const webW = Math.max(0.6, zoneWidth * 0.35 * lipScale);
  const lipW = Math.max(0.5, zoneWidth * 0.55 * lipScale);

  const withZ = (panel: Omit<RunnerPanelSpec, 'z'> & { z?: number }): RunnerPanelSpec => ({
    ...panel,
    z: runnerZ,
  });

  if (side === 'left') {
    return [
      withZ({
        x: zoneStartX,
        y: runnerY + lipH,
        width: webW,
        height: Math.max(0, runnerHeight - 2 * lipH),
        depth: runnerDepth,
        role,
      }),
      withZ({
        x: zoneStartX,
        y: runnerY,
        width: zoneWidth,
        height: lipH,
        depth: runnerDepth,
        role,
      }),
      withZ({
        x: zoneStartX + zoneWidth - lipW,
        y: runnerY + runnerHeight - lipH,
        width: lipW,
        height: lipH,
        depth: runnerDepth,
        role,
      }),
    ];
  }

  const zoneEnd = zoneStartX + zoneWidth;
  return [
    withZ({
      x: zoneEnd - webW,
      y: runnerY + lipH,
      width: webW,
      height: Math.max(0, runnerHeight - 2 * lipH),
      depth: runnerDepth,
      role,
    }),
    withZ({
      x: zoneStartX,
      y: runnerY,
      width: zoneWidth,
      height: lipH,
      depth: runnerDepth,
      role,
    }),
    withZ({
      x: zoneStartX,
      y: runnerY + runnerHeight - lipH,
      width: lipW,
      height: lipH,
      depth: runnerDepth,
      role,
    }),
  ];
}

/**
 * Corredera de carcasa (U grande) fija al interior del volumen.
 * Ocupa la zona [zoneStartX, zoneStartX + stackWidthMm].
 */
export function buildCabinetRunnerU(
  side: RunnerSide,
  zoneStartX: number,
  stackWidthMm: number,
  runnerY: number,
  runnerHeight: number,
  runnerDepth: number,
  runnerZ: number,
): RunnerPanelSpec[] {
  return uChannelPanels(
    zoneStartX,
    stackWidthMm,
    runnerY,
    runnerHeight,
    runnerDepth,
    runnerZ,
    'runner_cabinet',
    side,
    1,
  );
}

/**
 * Corredera de cajón (U chica), pegada al lateral exterior del cajón.
 */
export function buildDrawerRunnerU(
  side: RunnerSide,
  zoneStartX: number,
  stackWidthMm: number,
  runnerY: number,
  runnerHeight: number,
  runnerDepth: number,
  runnerZ: number,
): RunnerPanelSpec[] {
  return uChannelPanels(
    zoneStartX,
    stackWidthMm,
    runnerY,
    runnerHeight,
    runnerDepth,
    runnerZ,
    'runner_drawer',
    side,
    0.75,
  );
}

/** Zona corredera izquierda: inmediatamente a la izquierda del cajón. */
export function runnerZoneLeftX(boxX: number, stackWidthMm: number): number {
  return boxX - stackWidthMm;
}

/** Zona corredera derecha: borde derecho exterior del cajón. */
export function runnerZoneRightX(boxEndX: number): number {
  return boxEndX;
}
