import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import Scene3D from '../components/Scene3D';
import CustomCursor from '../components/CustomCursor';
import { Link } from 'react-router-dom';

gsap.registerPlugin(ScrollTrigger);

const FEATURES = [
  { title: 'Property Management', desc: 'Manage buildings, floors, zones, rooms, and beds dynamically.' },
  { title: 'Tenant & Identity', desc: 'Securely manage tenant identities, documents, and history.' },
  { title: 'Contracts & Renewals', desc: 'Flexible configurations for terms, deposits, and booking fees.' },
  { title: 'Financial Engine', desc: 'Automated invoices, payments, refunds, and charges.' },
  { title: 'Assets & Maintenance', desc: 'Track work orders, asset inspections, and vendor tasks.' },
  { title: 'Unified Dashboard', desc: 'Real-time analytics for occupancy, revenue, and operations.' },
];

export default function ImmersiveLanding() {
  const containerRef = useRef<HTMLDivElement>(null);
  const scrollProgress = useRef(0);

  useEffect(() => {
    // 1. Update 3D model progress via ScrollTrigger
    const st = ScrollTrigger.create({
      trigger: containerRef.current,
      start: 'top top',
      end: 'bottom bottom',
      onUpdate: (self) => {
        scrollProgress.current = self.progress;
      },
    });

    // 2. Fade-in animations for text
    const texts = gsap.utils.toArray('.animate-text');
    texts.forEach((text: any) => {
      gsap.fromTo(text, 
        { y: 80, opacity: 0 }, 
        { 
          y: 0, 
          opacity: 1,
          duration: 1.2,
          ease: 'power3.out',
          scrollTrigger: {
            trigger: text,
            start: 'top 85%',
          }
        }
      );
    });

    // 3. Staggered feature cards
    gsap.fromTo('.feature-card',
      { y: 50, opacity: 0 },
      {
        y: 0,
        opacity: 1,
        duration: 0.8,
        stagger: 0.15,
        ease: 'power3.out',
        scrollTrigger: {
          trigger: '.features-grid',
          start: 'top 75%',
        }
      }
    );

    return () => {
      st.kill();
      ScrollTrigger.getAll().forEach(t => t.kill());
    };
  }, []);

  return (
    <div ref={containerRef} className="relative min-h-[500vh] bg-[#0b0b0c] text-white selection:bg-orange/30">
      <CustomCursor />
      
      {/* 3D WebGL Background */}
      <Scene3D scrollProgress={scrollProgress} />

      {/* HTML Content Overlay */}
      {/* We use pointer-events-none on the wrapper so you can drag the 3D city through the empty spaces! */}
      <div className="relative z-10 overflow-hidden pointer-events-none">
        
        {/* Section 1: Hero */}
        <section className="min-h-screen flex flex-col items-center justify-center px-10 text-center">
          <div className="bg-black/30 backdrop-blur-md p-10 rounded-3xl border border-white/10 shadow-2xl max-w-5xl mx-auto">
            <div className="flex items-center justify-center gap-3 text-xs tracking-[0.2em] uppercase text-orange mb-6">
              <span className="w-6 h-px bg-orange"></span>
              Enterprise Property Management
              <span className="w-6 h-px bg-orange"></span>
            </div>
            <h1 className="animate-text font-display text-5xl md:text-7xl lg:text-8xl text-white leading-tight drop-shadow-2xl">
              Manage <em className="italic font-light text-white/80">Any Property.</em>
            </h1>
            <p className="animate-text mt-8 text-xl text-white/70 max-w-2xl mx-auto leading-relaxed">
              From Boarding Houses and Apartments to Commercial Buildings. Everything should be configurable. Nothing should be hardcoded.
            </p>
          </div>
        </section>

        {/* Section 2: The Problem */}
        <section className="min-h-screen flex flex-col items-center justify-center px-10 text-center">
          <div className="max-w-4xl mx-auto bg-black/40 backdrop-blur-lg p-12 rounded-3xl border border-white/5">
            <h2 className="animate-text font-display text-4xl md:text-5xl text-white mb-8 drop-shadow-lg">
              The Industry Problem: <br/> <span className="text-orange">Hardcoded & Fragmented</span>
            </h2>
            <p className="animate-text text-xl md:text-2xl text-white/70 font-light leading-relaxed">
              Most software is built for one specific business—forcing you to change systems when you grow from a single Kost to a multi-organization franchise. EPMP eliminates fragmented data and hardcoded rules, letting the software adapt to your business.
            </p>
          </div>
        </section>

        {/* Section 3: Core Capabilities (Bento Grid) */}
        <section className="min-h-screen flex flex-col items-center justify-center px-10 max-w-7xl mx-auto py-32 pointer-events-auto">
          <div className="bg-black/60 backdrop-blur-xl p-10 md:p-16 rounded-[3rem] border border-white/10 shadow-[0_0_50px_rgba(0,0,0,0.5)] w-full">
            <h2 className="animate-text font-display text-4xl md:text-5xl text-white mb-16 text-center drop-shadow-md">
              Unified Core Capabilities
            </h2>
            <div className="features-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full">
              {FEATURES.map((feat, idx) => (
                <div key={idx} className="feature-card pointer-events-auto border border-white/10 p-8 rounded-2xl bg-white/5 backdrop-blur-md hover:-translate-y-3 hover:bg-white/10 hover:border-orange/50 hover:shadow-[0_10px_30px_rgba(255,103,17,0.15)] transition-all duration-300 cursor-default group">
                  <div className="w-12 h-12 rounded-full border border-orange/50 mb-6 flex items-center justify-center text-orange group-hover:bg-orange group-hover:text-black group-hover:border-orange transition-colors duration-300 text-lg font-bold">
                    {idx + 1}
                  </div>
                  <h3 className="text-2xl font-bold text-white mb-3">{feat.title}</h3>
                  <p className="text-white/60 leading-relaxed text-base">{feat.desc}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Section 4: Sketchfab Embed & Interactive 3D */}
        <section className="min-h-screen flex flex-col lg:flex-row items-center justify-center gap-12 px-10 max-w-7xl mx-auto py-32 pointer-events-auto">
          <div className="flex-1 text-left bg-black/40 backdrop-blur-lg p-10 rounded-3xl border border-white/10">
            <h2 className="animate-text font-display text-4xl md:text-5xl leading-tight text-white mb-6">
              Immersive <br/> <span className="text-orange drop-shadow-[0_0_15px_rgba(255,103,17,0.5)]">3D Showcases.</span>
            </h2>
            <p className="animate-text text-lg text-white/70 mb-8 leading-relaxed">
              Native WebGL support alongside third-party integrations like Sketchfab. Provide prospective tenants with interactive, realistic walkthroughs of your assets before they sign the contract.
            </p>
          </div>
          
          <div className="flex-1 w-full animate-text">
            <div className="sketchfab-embed-wrapper w-full rounded-2xl overflow-hidden shadow-[0_0_40px_rgba(0,0,0,0.8)] border border-white/20 aspect-video bg-black hover:scale-105 transition-transform duration-500 ease-out">
              <iframe 
                title="Wire Frame City 001" 
                className="w-full h-full"
                frameBorder="0" 
                allowFullScreen 
                allow="autoplay; fullscreen; xr-spatial-tracking" 
                execution-while-out-of-viewport="true" 
                execution-while-not-rendered="true" 
                web-share="true" 
                src="https://sketchfab.com/models/ffe587c35191442082ea9a417abada2d/embed"> 
              </iframe>
            </div>
            <p className="text-xs text-center mt-4 text-white/50">
              Powered by Sketchfab Enterprise Embeds.
            </p>
          </div>
        </section>

        {/* Section 5: CTA / End */}
        <section className="section-end min-h-screen flex flex-col items-center justify-center px-10 text-center">
          <div className="bg-black/50 backdrop-blur-xl p-16 rounded-[3rem] border border-orange/20 shadow-[0_0_60px_rgba(255,103,17,0.1)]">
            <h2 className="animate-text font-display text-5xl md:text-7xl lg:text-8xl mb-8 text-white drop-shadow-xl">
              Ready to <span className="text-orange font-bold">Scale?</span>
            </h2>
            <p className="animate-text text-xl md:text-2xl text-white/70 mb-12 max-w-3xl mx-auto leading-relaxed">
              The foundation for AI Analytics, Dynamic Pricing, IoT Smart Locks, and your ultimate enterprise operational dashboard.
            </p>
            <div className="animate-text pointer-events-auto">
              <Link to="/auth/signin" className="inline-block bg-orange text-black font-bold text-xl px-12 py-5 rounded-full hover:bg-white hover:scale-110 hover:shadow-[0_0_30px_rgba(255,255,255,0.5)] transition-all duration-300">
                Enter EPMP Workspace
              </Link>
            </div>
          </div>
        </section>

      </div>
    </div>
  );
}
