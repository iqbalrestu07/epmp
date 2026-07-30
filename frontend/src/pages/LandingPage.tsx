
import { Link } from "react-router-dom";
import { ArrowRight, CheckCircle2, Building2, Users, Wallet } from "lucide-react";

export default function LandingPage() {
  return (
    <div className="min-h-screen text-foreground selection:bg-orange-bright/30">
      {/* ── HEADER ── */}
      <header className="absolute top-0 left-0 right-0 z-50 flex items-center justify-between px-6 py-5 max-w-7xl mx-auto">
        <div className="font-display text-xl tracking-[0.2em] uppercase text-foreground hover:opacity-70 transition-opacity">
          EPMP
        </div>
        <nav className="hidden md:flex items-center gap-8">
          <a href="#features" className="text-sm tracking-widest uppercase text-muted-foreground hover:text-foreground transition-colors relative group">
            Features
            <span className="absolute left-0 right-full bottom-0 h-px bg-orange transition-all duration-300 group-hover:right-0"></span>
          </a>
          <a href="#pricing" className="text-sm tracking-widest uppercase text-muted-foreground hover:text-foreground transition-colors relative group">
            Pricing
            <span className="absolute left-0 right-full bottom-0 h-px bg-orange transition-all duration-300 group-hover:right-0"></span>
          </a>
          <Link
            to="/dashboard"
            className="text-xs tracking-widest uppercase px-5 py-2.5 rounded-full bg-glass-bg border border-glass-border backdrop-blur-md text-muted-foreground hover:text-foreground hover:border-white/20 transition-all"
          >
            Dashboard
          </Link>
        </nav>
      </header>

      {/* ── HERO ── */}
      <main>
        <section className="relative pt-32 pb-24 min-h-[90svh] flex flex-col justify-center max-w-7xl mx-auto px-6">
          <div className="flex flex-col max-w-2xl">
            <div className="flex items-center gap-3 text-xs tracking-[0.2em] uppercase text-orange mb-6">
              <span className="w-6 h-px bg-orange"></span>
              Enterprise Ready
            </div>
            
            <h1 className="font-display text-5xl md:text-7xl lg:text-8xl leading-[0.95] tracking-wide mb-6">
              Property Management <br className="hidden md:block" />
              <span className="text-muted-foreground italic font-light">Evolved.</span>
            </h1>
            
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed max-w-xl mb-10">
              The same high-end architecture used in top-tier real estate operations. Ready to scale, adapt, and maximize your portfolio's ROI.
            </p>
            
            <div className="flex flex-col sm:flex-row items-start gap-4">
              <Link to="/dashboard" className="btn-primary w-full sm:w-auto">
                Enter Dashboard <ArrowRight className="w-5 h-5" />
              </Link>
              <a href="#features" className="h-14 px-6 flex items-center justify-center text-muted-foreground hover:text-foreground border-b border-white/20 hover:border-white/50 transition-colors">
                Explore Features
              </a>
            </div>
            
            <p className="mt-6 text-sm text-muted-foreground flex items-center gap-2">
              <CheckCircle2 className="w-4 h-4 text-orange" />
              <strong>No credit card required</strong> for early access.
            </p>
          </div>

          {/* Hero Image/Video Mockup */}
          <div className="mt-20 md:mt-32 aspect-[16/9] glass p-2 md:p-4 rounded-[2rem]">
            <div className="w-full h-full bg-[#0a0a0c] rounded-2xl border border-glass-border shadow-[0_30px_80px_rgba(0,0,0,0.6)] overflow-hidden relative">
              <div className="absolute top-4 right-4 flex items-center gap-2 px-3 py-1.5 rounded-full bg-glass-bg border border-glass-border backdrop-blur-md text-[10px] tracking-widest uppercase text-foreground z-10">
                <span className="w-2 h-2 rounded-full bg-orange shadow-[0_0_10px_#ff6711] animate-pulse"></span>
                Live Demo
              </div>
              {/* Placeholder for dashboard screenshot */}
              <div className="w-full h-full bg-gradient-to-br from-card to-background flex items-center justify-center">
                <span className="font-display text-4xl text-white/10">EPMP Dashboard Preview</span>
              </div>
            </div>
          </div>
        </section>

        {/* ── FEATURES ── */}
        <section id="features" className="py-24 max-w-7xl mx-auto px-6">
          <div className="mb-16">
            <h2 className="font-display text-4xl md:text-5xl lg:text-6xl mb-4">
              Built for <em className="italic text-foreground">Scale.</em>
            </h2>
            <p className="text-muted-foreground text-lg max-w-xl">
              Manage thousands of units, track complex financials, and keep your tenants happy with our modular architecture.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Feature 1 */}
            <div className="glass p-8 group hover:-translate-y-1 transition-transform duration-500">
              <div className="w-12 h-12 rounded-2xl bg-orange/10 border border-orange/20 flex items-center justify-center text-orange mb-6 group-hover:scale-110 transition-transform">
                <Building2 className="w-6 h-6" />
              </div>
              <h3 className="text-xl font-medium mb-3 text-foreground">Property Portfolio</h3>
              <p className="text-muted-foreground leading-relaxed">
                Organize buildings, floors, and rooms with ease. Track occupancy rates and maintenance schedules in real-time.
              </p>
            </div>
            
            {/* Feature 2 */}
            <div className="glass p-8 group hover:-translate-y-1 transition-transform duration-500">
              <div className="w-12 h-12 rounded-2xl bg-orange/10 border border-orange/20 flex items-center justify-center text-orange mb-6 group-hover:scale-110 transition-transform">
                <Users className="w-6 h-6" />
              </div>
              <h3 className="text-xl font-medium mb-3 text-foreground">Tenant CRM</h3>
              <p className="text-muted-foreground leading-relaxed">
                Complete tenant lifecycle management. From background checks and lease signing to communications and renewals.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="glass p-8 group hover:-translate-y-1 transition-transform duration-500">
              <div className="w-12 h-12 rounded-2xl bg-orange/10 border border-orange/20 flex items-center justify-center text-orange mb-6 group-hover:scale-110 transition-transform">
                <Wallet className="w-6 h-6" />
              </div>
              <h3 className="text-xl font-medium mb-3 text-foreground">Financial Engine</h3>
              <p className="text-muted-foreground leading-relaxed">
                Automated invoicing, payment tracking, and high-level reporting. Built-in hooks for popular accounting software.
              </p>
            </div>
          </div>
        </section>

        {/* ── CTA ── */}
        <section className="py-24 max-w-7xl mx-auto px-6 text-center">
          <div className="glass-strong p-12 md:p-20 flex flex-col items-center">
            <h2 className="font-display text-4xl md:text-5xl lg:text-6xl mb-6">
              Ready to elevate your operations?
            </h2>
            <p className="text-muted-foreground text-lg max-w-xl mb-10">
              Join top-tier real estate operators who manage their portfolios with EPMP.
            </p>
            <Link to="/dashboard" className="btn-primary w-full sm:w-auto">
              Get Started Now <ArrowRight className="w-5 h-5" />
            </Link>
          </div>
        </section>
      </main>

      {/* ── FOOTER ── */}
      <footer className="border-t border-white/10 py-10 text-center relative z-10">
        <div className="font-display text-lg tracking-[0.2em] uppercase text-foreground mb-4">EPMP</div>
        <p className="text-[11px] tracking-widest uppercase text-muted-foreground">
          &copy; {new Date().getFullYear()} Enterprise Property Management Platform. All rights reserved.
        </p>
      </footer>
    </div>
  );
}
