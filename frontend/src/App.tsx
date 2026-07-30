import { Routes, Route } from "react-router-dom";
import MainLayout from "./layouts/MainLayout";
import DashboardPage from "./pages/DashboardPage";
import ImmersiveLanding from "./features/immersive-view/pages/ImmersiveLanding";

export default function App() {
  return (
    <Routes>
      {/* Public Routes */}
      <Route path="/" element={<ImmersiveLanding />} />
      
      {/* Protected/App Routes */}
      <Route element={<MainLayout />}>
        <Route path="/dashboard" element={<DashboardPage />} />
      </Route>
    </Routes>
  );
}
