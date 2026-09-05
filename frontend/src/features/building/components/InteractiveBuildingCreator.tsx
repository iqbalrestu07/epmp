import { useState, useRef, Suspense } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls, Environment, ContactShadows, useGLTF } from "@react-three/drei";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePropertys } from "../../property/hooks";
import { createBuildingSchema, type CreateBuildingFormData } from "../schema";
import { Building2, Layers, MapPin, CheckCircle2, MousePointerClick, RefreshCw } from "lucide-react";
import * as THREE from "three";

useGLTF.preload("/3d-models/buildings/building.glb");

interface Props {
  onSubmit: (data: CreateBuildingFormData) => void;
  isSubmitting: boolean;
}

function GLBBuildingAsset({ targetY }: { targetY: number }) {
  const { scene } = useGLTF("/3d-models/buildings/building.glb");
  const cloned = useRef<THREE.Group>();

  if (!cloned.current) {
    cloned.current = scene.clone(true);
    cloned.current.traverse((child) => {
      if (child instanceof THREE.Mesh) {
        child.castShadow = true;
        child.receiveShadow = true;
      }
    });
  }

  const groupRef = useRef<THREE.Group>(null);

  useFrame((_, delta) => {
    if (groupRef.current) {
      groupRef.current.position.y = THREE.MathUtils.damp(
        groupRef.current.position.y,
        targetY,
        8,
        delta
      );
    }
  });

  return (
    <group ref={groupRef} position={[0, 15, 0]}>
      <primitive object={cloned.current} scale={[0.6, 0.6, 0.6]} />
    </group>
  );
}

export function InteractiveBuildingCreator({ onSubmit, isSubmitting }: Props) {
  const [buildingPlaced, setBuildingPlaced] = useState(false);
  const [position, setPosition] = useState<[number, number, number]>([0, 0, 0]);

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateBuildingFormData>({
    resolver: zodResolver(createBuildingSchema),
    defaultValues: {
      total_floors: 3,
      property_id: properties[0]?.id ?? "",
    },
  });


  const handleGroundClick = (e: any) => {
    e.stopPropagation();
    const snapX = Math.round(e.point.x);
    const snapZ = Math.round(e.point.z);
    setPosition([snapX, 0, snapZ]);
    setBuildingPlaced(true);
  };

  return (
    <div className="flex flex-col lg:flex-row gap-6 w-full h-[680px] bg-slate-100/60 p-4 rounded-2xl border border-slate-200">
      {/* 3D Canvas Area */}
      <div className="flex-1 relative bg-white rounded-xl overflow-hidden border border-slate-200 shadow-inner">
        {/* Helper instructions banner on top of canvas */}
        <div className="absolute top-4 left-4 right-4 z-10 flex items-center justify-between pointer-events-none">
          <div className="bg-white/95 backdrop-blur px-4 py-2.5 rounded-xl shadow-md border border-slate-200 pointer-events-auto flex items-center gap-3">
            {buildingPlaced ? (
              <>
                <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center text-green-600 shrink-0">
                  <CheckCircle2 size={18} />
                </div>
                <div>
                  <p className="text-xs font-semibold text-slate-800">Building Dropped at ({position[0]}, {position[2]})</p>
                  <p className="text-[11px] text-slate-500">Click anywhere on the ground to relocate, or fill attributes on the right.</p>
                </div>
              </>
            ) : (
              <>
                <div className="w-8 h-8 rounded-full bg-orange/10 flex items-center justify-center text-orange shrink-0 animate-pulse">
                  <MousePointerClick size={18} />
                </div>
                <div>
                  <p className="text-xs font-semibold text-slate-800">Step 1: Place 3D Building</p>
                  <p className="text-[11px] text-slate-500">Click anywhere on the ground grid to drop the building model.</p>
                </div>
              </>
            )}
          </div>

          {buildingPlaced && (
            <button
              onClick={() => setBuildingPlaced(false)}
              className="bg-white/95 backdrop-blur px-3 py-1.5 rounded-lg shadow-md border border-slate-200 text-xs font-medium text-slate-600 hover:text-slate-900 pointer-events-auto flex items-center gap-1.5 transition-colors"
            >
              <RefreshCw size={12} />
              Reset
            </button>
          )}
        </div>

        <Canvas shadows camera={{ position: [14, 12, 14], fov: 45 }}>
          <color attach="background" args={["#f8fafc"]} />
          <ambientLight intensity={0.6} />
          <directionalLight
            castShadow
            position={[15, 25, 10]}
            intensity={1.5}
            shadow-mapSize={[2048, 2048]}
            shadow-bias={-0.0001}
          />
          <Environment preset="city" />

          {/* Interactive Ground Plane */}
          <mesh
            rotation={[-Math.PI / 2, 0, 0]}
            position={[0, 0, 0]}
            receiveShadow
            onClick={handleGroundClick}
            onPointerOver={() => (document.body.style.cursor = "crosshair")}
            onPointerOut={() => (document.body.style.cursor = "default")}
          >
            <planeGeometry args={[40, 40]} />
            <meshStandardMaterial color="#e2e8f0" roughness={0.8} />
          </mesh>

          <gridHelper args={[40, 40, "#94a3b8", "#cbd5e1"]} position={[0, 0.01, 0]} />

          {/* Render 3D Model when placed */}
          {buildingPlaced && (
            <group position={position}>
              <Suspense
                fallback={
                  <mesh position={[0, 1.5, 0]} castShadow>
                    <boxGeometry args={[3, 3, 3]} />
                    <meshStandardMaterial color="#f97316" />
                  </mesh>
                }
              >
                <GLBBuildingAsset targetY={0} />
              </Suspense>

              {/* Glowing base ring */}
              <mesh position={[0, 0.03, 0]} rotation={[-Math.PI / 2, 0, 0]}>
                <ringGeometry args={[2.0, 2.3, 32]} />
                <meshBasicMaterial color="#f97316" side={THREE.DoubleSide} />
              </mesh>
            </group>
          )}

          <OrbitControls
            enableDamping
            dampingFactor={0.05}
            maxPolarAngle={Math.PI / 2 - 0.05}
            minDistance={4}
            maxDistance={35}
          />
          <ContactShadows position={[0, 0.02, 0]} opacity={0.4} scale={40} blur={2} far={10} />
        </Canvas>
      </div>

      {/* Docked Form Panel (Right Sidebar) */}
      <div className="w-full lg:w-96 bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 pb-4 border-b border-slate-100 mb-5">
            <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center text-orange">
              <Building2 size={18} />
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-base">Building Attributes</h3>
              <p className="text-xs text-slate-500">Step 2: Define properties and save</p>
            </div>
          </div>

          <form id="building-creator-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {/* Property Selector */}
            <div className="space-y-1.5">
              <Label htmlFor="property_id" className="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
                <MapPin size={13} className="text-orange" />
                Target Property
              </Label>
              <select
                id="property_id"
                {...register("property_id")}
                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-orange/30"
              >
                <option value="">Select a property…</option>
                {properties.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>
              {errors.property_id && (
                <p className="text-xs text-red-500">{errors.property_id.message}</p>
              )}
            </div>

            {/* Building Name */}
            <div className="space-y-1.5">
              <Label htmlFor="name" className="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
                <Building2 size={13} className="text-orange" />
                Building Name
              </Label>
              <Input
                id="name"
                {...register("name")}
                placeholder="e.g. Tower Alpha, Building 1"
                className="rounded-lg border-slate-300"
              />
              {errors.name && (
                <p className="text-xs text-red-500">{errors.name.message}</p>
              )}
            </div>

            {/* Total Floors */}
            <div className="space-y-1.5">
              <Label htmlFor="total_floors" className="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
                <Layers size={13} className="text-orange" />
                Total Floors
              </Label>
              <Input
                id="total_floors"
                type="number"
                min={1}
                max={150}
                {...register("total_floors", { valueAsNumber: true })}
                className="rounded-lg border-slate-300"
              />
              {errors.total_floors && (
                <p className="text-xs text-red-500">{errors.total_floors.message}</p>
              )}
            </div>

            {/* Status overview badge */}
            <div className="p-3 bg-slate-50 rounded-lg border border-slate-200 text-xs space-y-1 mt-4">
              <div className="flex justify-between">
                <span className="text-slate-500">3D Position:</span>
                <span className="font-mono font-medium text-slate-800">
                  {buildingPlaced ? `[${position[0]}, ${position[2]}]` : "Not placed yet"}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">Model:</span>
                <span className="font-medium text-slate-800">building.glb</span>
              </div>
            </div>
          </form>
        </div>

        {/* Submit action button */}
        <div className="pt-4 border-t border-slate-100">
          <Button
            type="submit"
            form="building-creator-form"
            disabled={isSubmitting || !buildingPlaced}
            className={`w-full py-2.5 font-semibold text-white rounded-xl shadow-md transition-all ${
              buildingPlaced
                ? "bg-orange hover:bg-orange/90 shadow-orange/20"
                : "bg-slate-300 cursor-not-allowed text-slate-500"
            }`}
          >
            {isSubmitting
              ? "Saving Building..."
              : buildingPlaced
              ? "Save Building"
              : "Drop Building to Enable Save"}
          </Button>
        </div>
      </div>
    </div>
  );
}
