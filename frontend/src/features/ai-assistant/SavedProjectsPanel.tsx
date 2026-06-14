import type { SavedProject } from '@/domain/project';
import { useProjectsStore } from '@/store/projects-store';
import './ai-assistant.css';

interface SavedProjectsPanelProps {
  activeProjectId: string | null;
  onLoad: (project: SavedProject) => void;
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('es', {
    dateStyle: 'short',
    timeStyle: 'short',
  });
}

export function SavedProjectsPanel({ activeProjectId, onLoad }: SavedProjectsPanelProps) {
  const projects = useProjectsStore((state) => state.projects);
  const deleteProject = useProjectsStore((state) => state.deleteProject);

  if (projects.length === 0) {
    return (
      <details className="saved-projects">
        <summary>Proyectos guardados</summary>
        <p className="saved-projects__empty">Aún no hay proyectos. Genera un diseño y guárdalo.</p>
      </details>
    );
  }

  return (
    <details className="saved-projects" open>
      <summary>Proyectos guardados ({projects.length})</summary>
      <ul className="saved-projects__list">
        {projects.map((project) => (
          <li key={project.id} className={project.id === activeProjectId ? 'is-active' : undefined}>
            <div className="saved-projects__meta">
              <strong>{project.name}</strong>
              <span>{formatDate(project.updatedAt)}</span>
              <p>{project.description}</p>
            </div>
            <div className="saved-projects__actions">
              <button type="button" className="btn btn--secondary" onClick={() => onLoad(project)}>
                Abrir
              </button>
              <button
                type="button"
                className="saved-projects__delete"
                onClick={() => deleteProject(project.id)}
              >
                Eliminar
              </button>
            </div>
          </li>
        ))}
      </ul>
    </details>
  );
}
