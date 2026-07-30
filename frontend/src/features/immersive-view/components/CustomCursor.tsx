import { useEffect, useRef } from 'react';
import gsap from 'gsap';

export default function CustomCursor() {
  const cursorRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Check if the device has a mouse
    if (window.matchMedia('(pointer: coarse)').matches) return;

    const xTo = gsap.quickTo(cursorRef.current, 'x', { duration: 0.15, ease: 'power3' });
    const yTo = gsap.quickTo(cursorRef.current, 'y', { duration: 0.15, ease: 'power3' });

    const moveCursor = (e: MouseEvent) => {
      xTo(e.clientX);
      yTo(e.clientY);
    };

    window.addEventListener('mousemove', moveCursor);
    return () => window.removeEventListener('mousemove', moveCursor);
  }, []);

  return (
    <div 
      ref={cursorRef} 
      className="fixed top-0 left-0 w-8 h-8 rounded-full border border-orange pointer-events-none z-50 mix-blend-difference -translate-x-1/2 -translate-y-1/2 transition-transform duration-300 hidden md:block"
    />
  );
}
