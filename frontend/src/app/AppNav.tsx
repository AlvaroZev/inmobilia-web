import { Link, useLocation } from 'react-router-dom';

const links = [
  { to: '/', label: 'Inicio' },
  { to: '/assistant', label: 'Asistente IA' },
  { to: '/viewer', label: 'Viewer 3D' },
];

export function AppNav() {
  const { pathname } = useLocation();

  return (
    <header className="app-nav">
      <Link to="/" className="app-nav__brand">
        Inmobilia
      </Link>
      <nav className="app-nav__links" aria-label="Principal">
        {links.map(({ to, label }) => (
          <Link
            key={to}
            to={to}
            className={pathname === to ? 'active' : undefined}
          >
            {label}
          </Link>
        ))}
      </nav>
    </header>
  );
}
