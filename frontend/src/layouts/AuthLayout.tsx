import { Outlet, Link } from 'react-router-dom';

export default function AuthLayout() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-[#0b0b0c] text-white selection:bg-orange/30 relative overflow-hidden">
      {/* Background decorations */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-orange/20 rounded-full blur-[120px] pointer-events-none"></div>
      <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-blue-500/10 rounded-full blur-[120px] pointer-events-none"></div>
      
      <div className="z-10 w-full max-w-md p-6">
        <div className="text-center mb-10">
          <Link to="/" className="inline-flex items-center gap-2 text-xs tracking-[0.2em] uppercase text-orange hover:text-white transition-colors">
            <span className="w-4 h-px bg-current"></span>
            Back to EPMP
            <span className="w-4 h-px bg-current"></span>
          </Link>
        </div>
        
        <div className="bg-black/40 backdrop-blur-xl border border-white/10 p-8 md:p-10 rounded-3xl shadow-2xl">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
