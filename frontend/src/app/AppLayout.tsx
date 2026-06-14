import type { ReactNode } from 'react';
import { AppNav } from './AppNav';
import './layout.css';

interface AppLayoutProps {
  children: ReactNode;
  full?: boolean;
}

export function AppLayout({ children, full }: AppLayoutProps) {
  return (
    <div className={`app-shell${full ? ' app-shell--full' : ''}`}>
      <AppNav />
      <main className="app-main">{children}</main>
    </div>
  );
}
