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

    // 4. Background color transition
    const bgSt = gsap.to(containerRef.current, {
      backgroundColor: '#f2efe9', 
      color: '#0b0b0c',
      scrollTrigger: {
        trigger: '.section-end',
        start: 'top center',
        end: 'bottom bottom',
        scrub: true,
      }
    });

    return () => {
      st.kill();
      bgSt.scrollTrigger?.kill();
      ScrollTrigger.getAll().forEach(t => t.kill());
    };
  }, []);

  return (
    <div ref={containerRef} className="relative min-h-[500vh] bg-background text-foreground transition-colors selection:bg-orange/30">
      <CustomCursor />
      
      {/* 3D WebGL Background */}
      <Scene3D scrollProgress={scrollProgress} />

      {/* HTML Content Overlay */}
      <div className="relative z-10 overflow-hidden">
        
        {/* Section 1: Hero */}
        <section className="h-screen flex flex-col items-center justify-center px-10 text-center">
          <div className="flex items-center gap-3 text-xs tracking-[0.2em] uppercase text-orange mb-6 mix-blend-difference">
            <span className="w-6 h-px bg-orange"></span>
            Enterprise Property Management
          </div>
          <h1 className="animate-text font-display text-5xl md:text-8xl lg:text-9xl mix-blend-difference pointer-events-none text-white leading-none">
            Manage <em className="italic font-light text-white/70">Any Property.</em>
          </h1>
          <p className="animate-text mt-8 text-lg text-white/60 max-w-xl mix-blend-difference">
            From Boarding Houses and Apartments to Commercial Buildings. Everything should be configurable. Nothing should be hardcoded.
          </p>
        </section>

        {/* Section 2: The Problem */}
        <section className="min-h-screen flex flex-col items-center justify-center px-10 text-center max-w-4xl mx-auto">
          <h2 className="animate-text font-display text-4xl md:text-5xl mix-blend-difference text-white mb-10">
            The Industry Problem: <br/> <span className="text-orange">Hardcoded & Fragmented</span>
          </h2>
          <p className="animate-text text-2xl text-white/70 font-light leading-relaxed mix-blend-difference">
            Most software is built for one specific business—forcing you to change systems when you grow from a single Kost to a multi-organization franchise. EPMP eliminates fragmented data and hardcoded rules, letting the software adapt to your business.
          </p>
        </section>

        {/* Section 3: Core Capabilities (Bento Grid) */}
        <section className="min-h-screen flex flex-col items-center justify-center px-10 max-w-7xl mx-auto py-32">
          <h2 className="animate-text font-display text-4xl md:text-5xl mix-blend-difference text-white mb-16 text-center">
            Unified Core Capabilities
          </h2>
          <div className="features-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 w-full mix-blend-difference">
            {FEATURES.map((feat, idx) => (
              <div key={idx} className="feature-card border border-white/20 p-8 rounded-2xl bg-black/20 backdrop-blur-md hover:bg-white/10 transition-colors cursor-default group">
                <div className="w-10 h-10 rounded-full border border-orange mb-6 flex items-center justify-center text-orange group-hover:bg-orange group-hover:text-black transition-colors">
                  {idx + 1}
                </div>
                <h3 className="text-xl font-bold text-white mb-3">{feat.title}</h3>
                <p className="text-white/60 leading-relaxed">{feat.desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Section 4: Sketchfab Embed & Interactive 3D */}
        <section className="min-h-screen flex flex-col lg:flex-row items-center justify-center gap-12 px-10 max-w-7xl mx-auto py-32">
          <div className="flex-1 text-left">
            <h2 className="animate-text font-display text-4xl md:text-5xl leading-tight mix-blend-difference text-white mb-6">
              Immersive <br/> <span className="text-orange">3D Showcases.</span>
            </h2>
            <p className="animate-text text-lg text-white/70 mix-blend-difference mb-8">
              Native WebGL support alongside third-party integrations like Sketchfab. Provide prospective tenants with interactive, realistic walkthroughs of your assets before they sign the contract.
            </p>
          </div>
          
          <div className="flex-1 w-full animate-text">
            {/* Sketchfab Embed */}
            <div className="sketchfab-embed-wrapper w-full rounded-2xl overflow-hidden shadow-2xl border border-white/20 aspect-video bg-black/50 backdrop-blur-md hover:scale-105 transition-transform duration-700 ease-out">
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
            <p className="text-xs text-center mt-4 text-white/40 mix-blend-difference">
              Powered by Sketchfab Enterprise Embeds.
            </p>
          </div>
        </section>

        {/* Section 5: CTA / End */}
        <section className="section-end min-h-screen flex flex-col items-center justify-center px-10 text-center">
          <h2 className="animate-text font-display text-6xl md:text-7xl lg:text-8xl mb-8 mix-blend-difference">
            Ready to <span className="text-orange font-bold">Scale?</span>
          </h2>
          <p className="animate-text text-2xl text-white/60 mb-12 max-w-3xl mix-blend-difference">
            The foundation for AI Analytics, Dynamic Pricing, IoT Smart Locks, and your ultimate enterprise operational dashboard.
          </p>
          <Link to="/dashboard" className="animate-text btn-primary text-black bg-white hover:bg-orange hover:text-white text-xl px-10 py-5 rounded-full transition-all duration-300">
            Enter EPMP Workspace
          </Link>
        </section>

      </div>
    </div>
  );
}
