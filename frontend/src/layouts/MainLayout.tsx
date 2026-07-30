import { Outlet } from "react-router-dom";

export default function MainLayout() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b bg-white px-6 py-4">
        <h1 className="text-lg font-semibold">EPMP</h1>
      </header>
      <main className="p-6">
        <Outlet />
      </main>
    </div>
  );
}
