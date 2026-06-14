import type { InstallationConstraints } from '@/domain/installation-constraints';
import type { RoomGeometry } from '@/domain/room-geometry';

export interface RoomSetupValues {
  roomWidth: number;
  roomDepth: number;
  roomHeight: number;
  nicheWidth: number;
  nicheHeight: number;
  nicheDepth: number;
  floorOffset: number;
  clearanceTop: number;
  clearanceLeft: number;
  clearanceRight: number;
  clearanceBack: number;
  clearanceFront: number;
}

const DEFAULT_TOLERANCES = { width: 2, height: 2, depth: 1 } as const;
const ANCHOR_WALL_ID = 'wall-north';

export const defaultRoomSetupValues = (): RoomSetupValues => ({
  roomWidth: 4200,
  roomDepth: 3600,
  roomHeight: 2700,
  nicheWidth: 4200,
  nicheHeight: 2700,
  nicheDepth: 620,
  floorOffset: 100,
  clearanceTop: 50,
  clearanceLeft: 50,
  clearanceRight: 50,
  clearanceBack: 10,
  clearanceFront: 0,
});

export function validateRoomSetup(values: RoomSetupValues): string[] {
  const errors: string[] = [];
  const positiveFields: Array<[keyof RoomSetupValues, string]> = [
    ['roomWidth', 'Ancho de habitación'],
    ['roomDepth', 'Profundidad de habitación'],
    ['roomHeight', 'Alto de habitación'],
    ['nicheWidth', 'Ancho del nicho'],
    ['nicheHeight', 'Alto del nicho'],
    ['nicheDepth', 'Profundidad del nicho'],
    ['floorOffset', 'Zócalo'],
    ['clearanceTop', 'Holgura superior'],
    ['clearanceLeft', 'Holgura izquierda'],
    ['clearanceRight', 'Holgura derecha'],
    ['clearanceBack', 'Holgura trasera'],
    ['clearanceFront', 'Holgura frontal'],
  ];

  for (const [field, label] of positiveFields) {
    if (!Number.isFinite(values[field]) || values[field] < 0) {
      errors.push(`${label} debe ser un número mayor o igual a 0.`);
    }
  }

  if (values.nicheWidth > values.roomWidth) {
    errors.push('El ancho del nicho no puede superar el ancho de la habitación.');
  }
  if (values.nicheHeight > values.roomHeight) {
    errors.push('El alto del nicho no puede superar el alto de la habitación.');
  }
  if (values.nicheDepth > values.roomDepth) {
    errors.push('La profundidad del nicho no puede superar la profundidad de la habitación.');
  }

  const effective = computeEffectiveInstallSpace(values);
  if (effective.width <= 0 || effective.height <= 0 || effective.depth <= 0) {
    errors.push('El espacio útil queda en cero o negativo. Revisa holguras y zócalo.');
  }

  return errors;
}

export function computeEffectiveInstallSpace(values: RoomSetupValues) {
  return {
    width:
      values.nicheWidth -
      values.clearanceLeft -
      values.clearanceRight -
      DEFAULT_TOLERANCES.width,
    height:
      values.nicheHeight -
      values.floorOffset -
      values.clearanceTop -
      DEFAULT_TOLERANCES.height,
    depth:
      values.nicheDepth -
      values.clearanceBack -
      values.clearanceFront -
      DEFAULT_TOLERANCES.depth,
  };
}

export function buildRoomGeometry(values: RoomSetupValues): RoomGeometry {
  const { roomWidth, roomDepth, roomHeight } = values;

  return {
    id: 'room-custom',
    name: 'Habitación personalizada',
    perimeter: {
      vertices: [
        { x: 0, y: 0 },
        { x: roomWidth, y: 0 },
        { x: roomWidth, y: roomDepth },
        { x: 0, y: roomDepth },
      ],
    },
    floor: {
      point: { x: 0, y: 0, z: 0 },
      normal: { x: 0, y: 1, z: 0 },
    },
    ceiling: {
      point: { x: 0, y: roomHeight, z: 0 },
      normal: { x: 0, y: -1, z: 0 },
    },
    walls: [
      {
        id: ANCHOR_WALL_ID,
        vertices: [
          { x: 0, y: 0, z: 0 },
          { x: roomWidth, y: 0, z: 0 },
          { x: roomWidth, y: roomHeight, z: 0 },
          { x: 0, y: roomHeight, z: 0 },
        ],
        thickness: 150,
      },
      {
        id: 'wall-east',
        vertices: [
          { x: roomWidth, y: 0, z: 0 },
          { x: roomWidth, y: 0, z: roomDepth },
          { x: roomWidth, y: roomHeight, z: roomDepth },
          { x: roomWidth, y: roomHeight, z: 0 },
        ],
        thickness: 150,
      },
    ],
    openings: [],
    obstacles: [
      {
        id: 'skirting-north',
        type: 'skirting',
        label: 'Zócalo pared norte',
        bounds: {
          min: { x: 0, y: 0, z: 0 },
          max: { x: values.nicheWidth, y: values.floorOffset, z: 50 },
        },
      },
    ],
  };
}

export function buildInstallationConstraints(values: RoomSetupValues): InstallationConstraints {
  return {
    id: 'install-custom',
    zone: {
      anchorWallIds: [ANCHOR_WALL_ID],
      bounds: {
        min: { x: 0, y: 0, z: 0 },
        max: { x: values.nicheWidth, y: values.nicheHeight, z: values.nicheDepth },
      },
    },
    clearances: {
      top: values.clearanceTop,
      bottom: 0,
      left: values.clearanceLeft,
      right: values.clearanceRight,
      back: values.clearanceBack,
      front: values.clearanceFront,
    },
    tolerances: { ...DEFAULT_TOLERANCES },
    references: {
      floorOffset: values.floorOffset,
      referenceWallId: ANCHOR_WALL_ID,
    },
  };
}

export function buildRoomAndInstallation(values: RoomSetupValues) {
  return {
    room: buildRoomGeometry(values),
    installation: buildInstallationConstraints(values),
  };
}
