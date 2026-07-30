# Immersive 3D Interactive Landing Page Guide
> Referensi Khusus untuk Modul 3D (Room/Building Viewer) di Proyek EPMP

Dokumen ini merangkum arsitektur, *tech stack*, dan kode dasar untuk membuat halaman *landing page* atau *viewer* interaktif 3D yang sangat imersif. Estetika yang dikejar adalah gabungan dari:
- **Ciao Energy**: Animasi *scrollytelling* yang halus dan integrasi teks kinetik.
- **Lusion.co**: Interaksi *mouse* berbasis *physics* dan efek 3D WebGL kelas atas.

Nantinya, konsep dan kode di dokumen ini dapat langsung diterapkan untuk membuat **Interactive Building/Room Viewer** bagi pengguna EPMP.

---

## 1. Tech Stack & Dependencies

Untuk mencapai tingkat imersi ini, Anda membutuhkan ekosistem **React Three Fiber (R3F)** digabungkan dengan **GSAP** untuk kontrol animasi *scroll* yang presisi.

```bash
# Jalankan perintah ini di dalam folder frontend/
npm install three @react-three/fiber @react-three/drei
npm install gsap @gsap/react
```

- **Three.js** & **@react-three/fiber**: *Engine* 3D berbasis komponen React.
- **@react-three/drei**: Kumpulan *helper* siap pakai (kamera, environment, kontrol) untuk R3F.
- **GSAP (GreenSock)**: Standar industri untuk animasi UI kompleks dan *scroll-trigger*.
- **Tailwind CSS**: Untuk tata letak HTML yang membungkus kanvas 3D.

---

## 2. Struktur Komponen (Modular)

Untuk menjaga agar aplikasi tetap performa tinggi, kita harus memisahkan komponen DOM (HTML) dan komponen Kanvas (WebGL).

```text
src/
 └── features/
     └── immersive-view/
         ├── components/
         │   ├── CustomCursor.tsx     # Custom cursor yang membesar saat interaksi
         │   ├── Scene3D.tsx          # Komponen Canvas utama (R3F)
         │   ├── FloatingMesh.tsx     # Objek 3D (Torus/Sphere) dengan material kaca
         │   └── ScrollContent.tsx    # Lapisan HTML (Teks Kinetik) di atas kanvas
         └── pages/
             └── ImmersiveLanding.tsx # Entry point
```

---

## 3. Implementasi Kode

### A. Custom Cursor (Micro-interactions)
Kursor kustom yang mengikuti pergerakan *mouse* menggunakan GSAP `quickTo` untuk performa 60FPS.

```tsx
// src/features/immersive-view/components/CustomCursor.tsx
import { useEffect, useRef } from 'react';
import gsap from 'context/gsap'; // (setup GSAP)

export default function CustomCursor() {
  const cursorRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
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
      className="fixed top-0 left-0 w-8 h-8 rounded-full border border-orange pointer-events-none z-50 mix-blend-difference -translate-x-1/2 -translate-y-1/2 transition-transform duration-300"
    />
  );
}
```

### B. Floating Mesh (Objek 3D & Mouse Tracking)
Membuat objek (contoh: *TorusKnot*) berbahan kaca yang secara halus merespon posisi kursor (Ala Lusion).

```tsx
// src/features/immersive-view/components/FloatingMesh.tsx
import { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import { MeshTransmissionMaterial } from '@react-three/drei';
import * as THREE from 'three';

export default function FloatingMesh({ scrollProgress }: { scrollProgress: React.MutableRefObject<number> }) {
  const meshRef = useRef<THREE.Mesh>(null);
  
  // Efek lerp mouse tracking
  const targetRotation = useRef({ x: 0, y: 0 });

  useFrame((state) => {
    if (!meshRef.current) return;

    // 1. Dapatkan kordinat mouse (-1 hingga 1)
    const { x, y } = state.pointer;
    
    // 2. Kalkulasi rotasi berdasar mouse (Lusion style)
    targetRotation.current.y = (x * Math.PI) / 4;
    targetRotation.current.x = (y * Math.PI) / 4;

    // 3. Smooth lerp ke rotasi target
    meshRef.current.rotation.x = THREE.MathUtils.lerp(meshRef.current.rotation.x, targetRotation.current.x, 0.05);
    meshRef.current.rotation.y = THREE.MathUtils.lerp(meshRef.current.rotation.y, targetRotation.current.y, 0.05);

    // 4. Scrollytelling efek (berputar dan membesar seiring scroll)
    const scale = 1 + scrollProgress.current * 2;
    meshRef.current.scale.set(scale, scale, scale);
    meshRef.current.position.y = -scrollProgress.current * 2;
  });

  return (
    <mesh ref={meshRef}>
      <torusKnotGeometry args={[1, 0.3, 128, 64]} />
      {/* Glassmorphism material */}
      <MeshTransmissionMaterial 
        backside
        thickness={1.5}
        roughness={0.1}
        transmission={1}
        ior={1.5}
        chromaticAberration={0.4}
        color="#ff6711"
      />
    </mesh>
  );
}
```

### C. Scene 3D Setup (R3F Canvas)
Menyiapkan Kanvas yang merender objek 3D di belakang konten HTML.

```tsx
// src/features/immersive-view/components/Scene3D.tsx
import { Canvas } from '@react-three/fiber';
import { Environment, Float } from '@react-three/drei';
import FloatingMesh from './FloatingMesh';

export default function Scene3D({ scrollProgress }: { scrollProgress: React.MutableRefObject<number> }) {
  return (
    <div className="fixed inset-0 z-0 pointer-events-none">
      <Canvas camera={{ position: [0, 0, 5], fov: 45 }}>
        <ambientLight intensity={0.5} />
        <directionalLight position={[10, 10, 5]} intensity={1} />
        
        {/* Floating effect tambahan */}
        <Float speed={2} rotationIntensity={1} floatIntensity={2}>
          <FloatingMesh scrollProgress={scrollProgress} />
        </Float>
        
        <Environment preset="city" />
      </Canvas>
    </div>
  );
}
```

### D. Scrollytelling HTML & Entry Point (GSAP)
Integrasi DOM dan 3D. GSAP ScrollTrigger mengontrol variabel `scrollProgress` yang diteruskan ke kanvas 3D, serta memicu animasi teks.

```tsx
// src/features/immersive-view/pages/ImmersiveLanding.tsx
import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import Scene3D from '../components/Scene3D';
import CustomCursor from '../components/CustomCursor';

gsap.registerPlugin(ScrollTrigger);

export default function ImmersiveLanding() {
  const containerRef = useRef<HTMLDivElement>(null);
  const scrollProgress = useRef(0);

  useEffect(() => {
    // 1. Update 3D model progress via ScrollTrigger
    ScrollTrigger.create({
      trigger: containerRef.current,
      start: 'top top',
      end: 'bottom bottom',
      onUpdate: (self) => {
        scrollProgress.current = self.progress;
      },
    });

    // 2. Animasi fade-in teks kinetik
    gsap.utils.toArray('.animate-text').forEach((text: any) => {
      gsap.fromTo(text, 
        { y: 100, opacity: 0 }, 
        { 
          y: 0, 
          opacity: 1,
          duration: 1,
          ease: 'power3.out',
          scrollTrigger: {
            trigger: text,
            start: 'top 80%',
          }
        }
      );
    });

    // 3. Background transisi warna (Ciao Energy style)
    gsap.to(containerRef.current, {
      backgroundColor: '#f2efe9', // Terang di akhir
      color: '#0b0b0c',
      scrollTrigger: {
        trigger: '.section-end',
        start: 'top center',
        end: 'bottom bottom',
        scrub: true,
      }
    });
  }, []);

  return (
    <div ref={containerRef} className="relative bg-background text-foreground transition-colors">
      <CustomCursor />
      <Scene3D scrollProgress={scrollProgress} />

      {/* Layer HTML di atas 3D */}
      <div className="relative z-10">
        
        {/* Section 1: Hero */}
        <section className="h-screen flex items-center justify-center px-10">
          <h1 className="animate-text font-display text-7xl md:text-9xl mix-blend-difference pointer-events-none text-white">
            Future Spaces.
          </h1>
        </section>

        {/* Section 2: Info */}
        <section className="h-screen flex items-center justify-start px-10">
          <p className="animate-text max-w-2xl text-2xl md:text-4xl font-light mix-blend-difference text-white">
            Interactive building visualizations that react to your presence.
          </p>
        </section>

        {/* Section 3: End / Change Color */}
        <section className="section-end h-screen flex items-center justify-end px-10">
          <h2 className="animate-text font-display text-6xl text-right">
            Explore <br/> Properties
          </h2>
        </section>

      </div>
    </div>
  );
}
```

---

## 4. Panduan Mengganti Asset 3D & Optimasi

Ketika desainer atau arsitek Anda memberikan file model 3D baru (misal: `gedung_baru.glb`), Anda **tidak perlu** menulis kode renderingnya secara manual! 

### Langkah 1: Eksekusi Script GLTF
Letakkan file mentah Anda di `frontend/public/models/`. Lalu jalankan perintah kompresi dan auto-generate yang sudah saya tanam di `package.json`:

```bash
npm run gltf public/models/gedung_baru.glb -o src/features/immersive-view/components/GedungBaru.tsx
```

**Apa yang terjadi di balik layar?**
- Alat ini akan mengkompresi model 3D menggunakan kompresi Draco (menekan ukuran file hingga 80-90%).
- Menghasilkan file terkompresi `gedung_baru-transformed.glb` secara otomatis.
- Menuliskan komponen React `GedungBaru.tsx` yang secara presisi me-load setiap `mesh` dan `material` di dalam model tersebut.

### Langkah 2: Gunakan di Scene Anda
Buka file `FloatingMesh.tsx` (atau komponen apapun tempat Anda ingin meletakkan model tersebut). 
Ganti *import* model lama dengan model baru:

```tsx
import { Model as GedungBaru } from './GedungBaru';

export default function FloatingMesh() {
  return (
    <group>
      {/* Panggil komponen hasil generate tadi */}
      <GedungBaru scale={1.5} position={[0, -2, 0]} />
    </group>
  );
}
```

### Langkah 3: Beri Efek Dramatis (Semenarik Mungkin)
Sebuah model 3D mentah tidak akan terlihat bagus tanpa tata cahaya dan efek post-processing lingkungan (Lingkungan / Environment). Untuk model *City Night* yang baru kita tambahkan, saya menggunakan efek ini di `Scene3D.tsx`:

- **Fog (Kabut)**: `<fog attach="fog" args={['#0b0b0c', 5, 20]} />` (menciptakan efek *fade-out* sinematik di ujung horizon, menyembunyikan batas tajam 3D).
- **SpotLight**: `<spotLight position={[0, 10, 0]} intensity={2} color="#ffaa00" />` (memberikan cahaya sorot bernuansa hangat ke pusat kota).
- **Float**: Menggunakan komponen `<Float>` dari *Drei* untuk membuat keseluruhan kota sedikit mengambang dan merespon gravitasi secara halus layaknya halusinasi.

Kombinasi antara interaksi *mouse tracking*, efek *Scrollytelling GSAP*, dan pencahayaan sinematik di atas lah yang membuat Landing Page ini terasa sangat mahal dan *Immersive* layaknya studio Lusion!

---

## 5. Alternatif Cepat: Menggunakan Iframe Embed (Sketchfab)

Jika Anda memiliki model 3D yang sangat berat atau Anda tidak ingin meng-_host_ file `.glb` tersebut di server sendiri, Anda bisa menggunakan layanan pihak ketiga seperti **Sketchfab**.

Sketchfab memberikan kemudahan dengan menyediakan tag `<iframe>` yang sudah mendukung *WebGL rendering* yang di-optimasi oleh mereka, termasuk *Post-Processing* dan *Lighting*.

**Cara Integrasi di React (Tailwind):**
1. Salin kode *embed* dari Sketchfab.
2. Letakkan di atas layer 3D R3F Anda (di dalam HTML biasa).
3. Bungkus dengan `<div>` styling seperti ini agar terlihat premium:

```tsx
<div className="sketchfab-embed-wrapper w-full rounded-2xl overflow-hidden shadow-2xl border border-white/20 aspect-video bg-black/50 backdrop-blur-md">
  <iframe 
    title="Wire Frame City 001" 
    className="w-full h-full"
    frameBorder="0" 
    allowFullScreen 
    src="https://sketchfab.com/models/ffe587c35191442082ea9a417abada2d/embed"> 
  </iframe>
</div>
```

**Kapan menggunakan R3F (WebGL murni) vs Iframe Embed?**
- Gunakan **R3F (Native WebGL)**: Jika Anda ingin model tersebut *react* (berinteraksi) dengan kursor mouse secara presisi, bereaksi terhadap *scroll* GSAP, atau jika Anda ingin menambahkan tombol interaktif langsung menempel di atas gedung (seperti *floor selector*).
- Gunakan **Iframe (Sketchfab)**: Untuk sekedar _showcase_ statis di mana pengguna hanya memutar-mutar model secara manual, tanpa perlu interaksi yang terhubung dengan state React Anda.
