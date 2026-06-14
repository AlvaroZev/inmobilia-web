import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { AppLayout } from './AppLayout';
import { AssistantPage } from '@/pages/AssistantPage';
import { HomePage } from '@/pages/HomePage';
import { ViewerPage } from '@/pages/ViewerPage';

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={
            <AppLayout>
              <HomePage />
            </AppLayout>
          }
        />
        <Route
          path="/assistant"
          element={
            <AppLayout>
              <AssistantPage />
            </AppLayout>
          }
        />
        <Route
          path="/viewer"
          element={
            <AppLayout full>
              <ViewerPage />
            </AppLayout>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
