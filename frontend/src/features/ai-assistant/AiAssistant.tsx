import { useState } from 'react';
import { Link } from 'react-router-dom';
import type { PipelineSnapshot, SavedProject } from '@/domain/project';
import { parseFurnitureDescription, solveFurniture, compileManufacturing, calculateCost } from '@/services/api';
import { ExportPanel } from '@/features/cost-estimator';
import { RoomSetupPanel } from '@/features/ai-assistant/RoomSetupPanel';
import { SavedProjectsPanel } from '@/features/ai-assistant/SavedProjectsPanel';
import { useAssistantStore } from '@/store/assistant-store';
import { useProjectsStore } from '@/store/projects-store';
import { useViewerStore } from '@/store/viewer-store';
import { normalizeFurnitureDefinition } from '@/utils/load-furniture';
import { buildRoomAndInstallation, validateRoomSetup } from '@/utils/room-setup';
import { validateFurnitureDefinition } from '@/utils/volume-tree';
import './ai-assistant.css';

export function AiAssistant() {
  const [loading, setLoading] = useState(false);

  const messages = useAssistantStore((s) => s.messages);
  const input = useAssistantStore((s) => s.input);
  const name = useAssistantStore((s) => s.name);
  const runPipeline = useAssistantStore((s) => s.runPipeline);
  const result = useAssistantStore((s) => s.result);
  const error = useAssistantStore((s) => s.error);
  const roomSetup = useAssistantStore((s) => s.roomSetup);
  const roomSetupErrors = useAssistantStore((s) => s.roomSetupErrors);
  const activeProjectId = useAssistantStore((s) => s.activeProjectId);
  const saveNotice = useAssistantStore((s) => s.saveNotice);

  const setInput = useAssistantStore((s) => s.setInput);
  const setName = useAssistantStore((s) => s.setName);
  const setRunPipeline = useAssistantStore((s) => s.setRunPipeline);
  const setResult = useAssistantStore((s) => s.setResult);
  const setError = useAssistantStore((s) => s.setError);
  const setRoomSetup = useAssistantStore((s) => s.setRoomSetup);
  const setRoomSetupErrors = useAssistantStore((s) => s.setRoomSetupErrors);
  const setActiveProjectId = useAssistantStore((s) => s.setActiveProjectId);
  const setSaveNotice = useAssistantStore((s) => s.setSaveNotice);
  const setMessages = useAssistantStore((s) => s.setMessages);
  const loadSession = useAssistantStore((s) => s.loadSession);
  const resetSession = useAssistantStore((s) => s.resetSession);

  const upsertProject = useProjectsStore((state) => state.upsertProject);
  const setResolvedFurniture = useViewerStore((state) => state.setResolvedFurniture);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!input.trim() || loading) return;

    const userMessage = {
      id: `user-${Date.now()}`,
      role: 'user' as const,
      content: input.trim(),
    };
    setMessages((prev) => [...prev, userMessage]);
    setLoading(true);
    setError(null);
    setSaveNotice(null);

    try {
      const raw = await parseFurnitureDescription({
        description: userMessage.content,
        name: name.trim() || undefined,
      });
      const furniture = normalizeFurnitureDefinition(raw);

      const validation = validateFurnitureDefinition(furniture);
      if (!validation.valid) {
        throw new Error('La IA devolvió un árbol de volúmenes inválido');
      }

      const pipeline: PipelineSnapshot = { furniture };

      if (runPipeline) {
        const setupErrors = validateRoomSetup(roomSetup);
        if (setupErrors.length > 0) {
          setRoomSetupErrors(setupErrors);
          throw new Error('Revisa las medidas del espacio de instalación.');
        }
        setRoomSetupErrors([]);

        const { room, installation } = buildRoomAndInstallation(roomSetup);
        pipeline.resolved = await solveFurniture({
          room,
          furniture,
          installation,
        });
        pipeline.manufacturing = await compileManufacturing(pipeline.resolved);
        pipeline.cost = await calculateCost(pipeline.manufacturing);
        setResolvedFurniture(pipeline.resolved);
      }

      setResult(pipeline);
      setMessages((prev) => [
        ...prev,
        {
          id: `assistant-${Date.now()}`,
          role: 'assistant',
          content: runPipeline
            ? `Listo: "${furniture.name}" con ${pipeline.manufacturing?.parts.length ?? 0} piezas. Costo total: ${pipeline.cost?.currency} ${pipeline.cost?.total.toFixed(2)}.`
            : `Listo: "${furniture.name}" — FurnitureDefinition validado.`,
        },
      ]);
      setInput('');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Error desconocido';
      setError(message);
      setMessages((prev) => [
        ...prev,
        { id: `error-${Date.now()}`, role: 'assistant', content: `Error: ${message}` },
      ]);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveProject = () => {
    if (!result) return;

    const projectName = name.trim() || result.furniture.name || 'Proyecto sin nombre';
    const id = upsertProject({
      id: activeProjectId ?? undefined,
      name: projectName,
      description: result.furniture.description ?? projectName,
      messages,
      roomSetup,
      runPipeline,
      result,
    });

    setActiveProjectId(id);
    setName(projectName);
    setSaveNotice('Proyecto guardado en este navegador.');
  };

  const handleLoadProject = (project: SavedProject) => {
    loadSession({
      activeProjectId: project.id,
      messages: project.messages,
      name: project.name,
      roomSetup: project.roomSetup,
      runPipeline: project.runPipeline,
      result: project.result,
    });
    setRoomSetupErrors(validateRoomSetup(project.roomSetup));

    if (project.result.resolved) {
      setResolvedFurniture(project.result.resolved);
    }
  };

  const handleNewProject = () => {
    resetSession();
  };

  return (
    <div className="ai-layout">
      <section className="ai-chat">
        <header>
          <p className="eyebrow">Asistente IA</p>
          <h1>Diseño por lenguaje natural</h1>
          <p className="ai-subtitle">Solo genera FurnitureDefinition — nunca geometría Three.js ni piezas directas.</p>
        </header>

        <SavedProjectsPanel activeProjectId={activeProjectId} onLoad={handleLoadProject} />

        <div className="ai-messages">
          {messages.map((message) => (
            <article key={message.id} className={`ai-message ai-message--${message.role}`}>
              <span>{message.role === 'user' ? 'Vendedor' : 'IA'}</span>
              <p>{message.content}</p>
            </article>
          ))}
        </div>

        <RoomSetupPanel
          values={roomSetup}
          onChange={(values) => {
            setRoomSetup(values);
            if (roomSetupErrors.length > 0) {
              setRoomSetupErrors(validateRoomSetup(values));
            }
          }}
          errors={roomSetupErrors}
        />

        <form className="ai-form" onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Nombre del proyecto (opcional)"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <textarea
            placeholder="Ej: Closet 2.4m dos cuerpos · Escritorio con cajonera · Centro de entretenimiento para TV..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            rows={4}
          />
          <label className="ai-checkbox">
            <input
              type="checkbox"
              checked={runPipeline}
              onChange={(e) => setRunPipeline(e.target.checked)}
            />
            Ejecutar solver + manufacturing + costo
          </label>
          <div className="ai-form__actions">
            <button type="submit" className="btn btn--primary" disabled={loading || !input.trim()}>
              {loading ? 'Procesando...' : 'Generar diseño'}
            </button>
            <button type="button" className="btn btn--secondary" onClick={handleNewProject}>
              Nuevo proyecto
            </button>
          </div>
        </form>

        {error && <p className="ai-error">{error}</p>}
      </section>

      <aside className="ai-output">
        <div className="ai-output__header">
          <h2>Salida</h2>
          {activeProjectId && <span className="ai-output__badge">Proyecto activo</span>}
        </div>
        {!result && <p className="ai-placeholder">El JSON aparecerá aquí después de generar.</p>}
        {result && (
          <>
            <div className="ai-viewer-actions">
              <button type="button" className="btn btn--primary" onClick={handleSaveProject}>
                {activeProjectId ? 'Actualizar proyecto' : 'Guardar proyecto'}
              </button>
              {result.resolved && (
                <Link to="/viewer" className="btn btn--secondary">
                  Ver en 3D
                </Link>
              )}
            </div>
            {saveNotice && <p className="ai-save-notice">{saveNotice}</p>}
            <details open>
              <summary>FurnitureDefinition</summary>
              <pre>{JSON.stringify(result.furniture, null, 2)}</pre>
            </details>
            {result.resolved && (
              <details>
                <summary>ResolvedFurniture</summary>
                <pre>{JSON.stringify(result.resolved, null, 2)}</pre>
              </details>
            )}
            {result.cost && (
              <div className="ai-cost-card">
                <strong>Costo total</strong>
                <p>
                  {result.cost.currency} {result.cost.total.toFixed(2)}
                </p>
                <small>
                  {result.manufacturing?.parts.length ?? 0} piezas · {result.cost.labor.hours} h mano de obra
                </small>
              </div>
            )}
            {result.manufacturing && (
              <ExportPanel
                furnitureName={result.furniture.name}
                model={result.manufacturing}
                cost={result.cost}
              />
            )}
          </>
        )}
      </aside>
    </div>
  );
}
