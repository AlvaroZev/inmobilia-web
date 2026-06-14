import {
  computeEffectiveInstallSpace,
  defaultRoomSetupValues,
  type RoomSetupValues,
} from '@/utils/room-setup';
import './ai-assistant.css';

interface RoomSetupPanelProps {
  values: RoomSetupValues;
  onChange: (values: RoomSetupValues) => void;
  errors?: string[];
}

type NumericField = keyof RoomSetupValues;

const FIELDS: Array<{ key: NumericField; label: string; hint?: string }> = [
  { key: 'roomWidth', label: 'Ancho habitación (mm)' },
  { key: 'roomDepth', label: 'Profundidad habitación (mm)' },
  { key: 'roomHeight', label: 'Alto habitación (mm)' },
  { key: 'nicheWidth', label: 'Ancho nicho (mm)' },
  { key: 'nicheHeight', label: 'Alto nicho (mm)' },
  { key: 'nicheDepth', label: 'Profundidad nicho (mm)' },
  { key: 'floorOffset', label: 'Zócalo (mm)' },
  { key: 'clearanceTop', label: 'Holgura superior (mm)' },
  { key: 'clearanceLeft', label: 'Holgura izquierda (mm)' },
  { key: 'clearanceRight', label: 'Holgura derecha (mm)' },
  { key: 'clearanceBack', label: 'Holgura trasera (mm)' },
  { key: 'clearanceFront', label: 'Holgura frontal (mm)' },
];

export function RoomSetupPanel({ values, onChange, errors = [] }: RoomSetupPanelProps) {
  const effective = computeEffectiveInstallSpace(values);

  const handleFieldChange = (key: NumericField, raw: string) => {
    const parsed = Number(raw);
    onChange({
      ...values,
      [key]: Number.isFinite(parsed) ? parsed : 0,
    });
  };

  const handleReset = () => {
    onChange(defaultRoomSetupValues());
  };

  return (
    <details className="room-setup" open>
      <summary>Espacio de instalación</summary>
      <p className="room-setup__intro">
        Define las medidas del ambiente y del nicho antes de ejecutar el solver.
      </p>

      <div className="room-setup__grid">
        {FIELDS.map(({ key, label }) => (
          <label key={key} className="room-setup__field">
            <span>{label}</span>
            <input
              type="number"
              min={0}
              step={10}
              value={values[key]}
              onChange={(e) => handleFieldChange(key, e.target.value)}
            />
          </label>
        ))}
      </div>

      <div className="room-setup__effective">
        <span>Espacio útil estimado</span>
        <strong>
          {Math.round(effective.width)} × {Math.round(effective.height)} × {Math.round(effective.depth)} mm
        </strong>
      </div>

      {errors.length > 0 && (
        <ul className="room-setup__errors">
          {errors.map((message) => (
            <li key={message}>{message}</li>
          ))}
        </ul>
      )}

      <button type="button" className="btn btn--secondary room-setup__reset" onClick={handleReset}>
        Restaurar valores demo
      </button>
    </details>
  );
}
