import { Link } from 'react-router-dom';

export function HomePage() {
  return (
    <section className="home">
      <p className="home__eyebrow">Plataforma paramétrica</p>
      <h1>Muebles de melamina a medida</h1>
      <p className="home__lead">
        Diseño por lenguaje natural, solver de restricciones, fabricación, costos y
        visualización 3D. El viewer renderiza exclusivamente{' '}
        <code>ResolvedFurniture</code>.
      </p>
      <div className="home__actions">
        <Link to="/assistant" className="btn btn--primary">
          Asistente IA
        </Link>
        <Link to="/viewer" className="btn btn--secondary">
          Viewer 3D
        </Link>
      </div>
    </section>
  );
}
