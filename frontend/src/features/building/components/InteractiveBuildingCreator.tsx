import { useState, useRef, Suspense } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls, Environment, ContactShadows, useGLTF, Html } from "@react-three/drei";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePropertys } from "../../property/hooks";
import { createBuildingSchema, type CreateBuildingFormData } from "../schema";
import { Building2, Layers, MapPin, Sparkles, Check, RefreshCw } from "lucide-react";
import * as THREE from "three";

// Preload both user-provided GLB models
useGLTF.preload("/3d-models/buildings/building.glb");
useGLTF.preload("/3d-models/buildings/hotel.glb");

export interface BuildingModelPreset {
  id: string;
  name: string;
  category: string;
  modelPath: string;
  scale: [number, number, number];
  previewIcon: string;
  tag: string;
  badgeColor: string;
  suggestedName: string;
}

export const BUILDING_PRESETS: BuildingModelPreset[] = [
  {
    id: "office_tower",
    name: "Modern Office Tower",
    category: "Commercial & Office",
    modelPath: "/3d-models/buildings/building.glb",
    scale: [0.6, 0.6, 0.6],
    previewIcon: "🏢",
    tag: "building.glb",
    badgeColor: "bg-blue-50 text-blue-700 border-blue-200",
    suggestedName: "Tower Alpha",
  },
  {
    id: "hotel_resort",
    name: "Hotel & Luxury Suites",
    category: "Hospitality & Resort",
    modelPath: "/3d-models/buildings/hotel.glb",
    scale: [0.55, 0.55, 0.55],
    previewIcon: "🏨",
    tag: "hotel.glb",
    badgeColor: "bg-amber-50 text-amber-800 border-amber-200",
    suggestedName: "Grand Hotel",
  },
];

interface Props {
  onSubmit: (data: CreateBuildingFormData) => void;
  isSubmitting: boolean;
}

// Dynamic 3D model asset renderer for chosen GLB
function GLBBuildingAsset({
  modelPath,
  scale,
  targetY,
}: {
  modelPath: string;
  scale: [number, number, number];
  targetY: number;
}) {
  const { scene } = useGLTF(modelPath);
  const cloned = scene.clone(true);

  cloned.traverse((child) => {
    if (child instanceof THREE.Mesh) {
      child.castShadow = true;
      child.receiveShadow = true;
    }
  });

  const groupRef = useRef<THREE.Group>(null);

  useFrame((_, delta) => {
    if (groupRef.current) {
      groupRef.current.position.y = THREE.MathUtils.damp(
        groupRef.current.position.y,
        targetY,
        10,
        delta
      );
    }
  });

  return (
    <group ref={groupRef} position={[0, 4, 0]}>
      <primitive object={cloned} scale={scale} position={[0, 0, 0]} />
    </group>
  );
}

export function InteractiveBuildingCreator({ onSubmit, isSubmitting }: Props) {
  const [selectedPreset, setSelectedPreset] = useState<BuildingModelPreset>(BUILDING_PRESETS[0]);
  const [position, setPosition] = useState<[number, number, number]>([0, 0, 0]);
  const [buildingPlaced, setBuildingPlaced] = useState(true);

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateBuildingFormData>({
    resolver: zodResolver(createBuildingSchema),
    defaultValues: {
      name: "Tower Alpha",
      total_floors: 4,
      property_id: properties[0]?.id ?? "",
    },
  });

  const currentBuildingName = watch("name") || selectedPreset.suggestedName;
  const currentFloors = watch("total_floors") || 4;

  const handleSelectPreset = (preset: BuildingModelPreset) => {
    setSelectedPreset(preset);
    const currentName = watch("name");
    if (!currentName || currentName === "Tower Alpha" || currentName === "Grand Hotel") {
      setValue("name", preset.suggestedName);
    }
  };

  const handleGroundClick = (e: any) => {
    e.stopPropagation();
    const snapX = Math.round(e.point.x);
    const snapZ = Math.round(e.point.z);
    setPosition([snapX, 0, snapZ]);
    setBuildingPlaced(true);
  };

  const handleFormSubmit = (data: CreateBuildingFormData) => {
    try {
      localStorage.setItem(`building_model_pref_${data.name}`, selectedPreset.id);
    } catch {
      // ignore
    }
    onSubmit(data);
  };

  return (
    <div className="flex flex-col lg:flex-row gap-6 w-full min-h-[720px] bg-slate-100/60 p-4 rounded-2xl border border-slate-200">
      {/* 3D Canvas Area */}
      <div className="flex-1 relative bg-white rounded-xl overflow-hidden border border-slate-200 shadow-inner min-h-[450px]">
        {/* Helper instructions banner on top of canvas */}
        <div className="absolute top-4 left-4 right-4 z-10 flex items-center justify-between pointer-events-none">
          <div className="bg-white/95 backdrop-blur px-4 py-2.5 rounded-xl shadow-md border border-slate-200 pointer-events-auto flex items-center gap-3">
            <div className="w-9 h-9 rounded-full bg-orange/10 flex items-center justify-center text-orange shrink-0">
              <span className="text-lg">{selectedPreset.previewIcon}</span>
            </div>
            <div>
              <p className="text-xs font-bold text-slate-800 flex items-center gap-1.5">
                {selectedPreset.name}
                <span className="text-[10px] px-1.5 py-0.5 rounded bg-orange/10 text-orange font-normal">
                  ({position[0]}, {position[2]})
                </span>
              </p>
              <p className="text-[11px] text-slate-500">
                Click anywhere on the ground grid to relocate building position.
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 pointer-events-auto">
            <button
              onClick={() => setPosition([0, 0, 0])}
              className="bg-white/95 backdrop-blur px-3 py-1.5 rounded-lg shadow-md border border-slate-200 text-xs font-medium text-slate-600 hover:text-slate-900 flex items-center gap-1.5 transition-colors"
            >
              <RefreshCw size={12} />
              Reset Center [0,0]
            </button>
          </div>
        </div>

        {/* Floating Model Badge indicator at bottom left */}
        <div className="absolute bottom-4 left-4 z-10 bg-white/90 backdrop-blur px-3 py-1.5 rounded-lg border shadow-sm flex items-center gap-2 text-xs">
          <span className="text-base">{selectedPreset.previewIcon}</span>
          <div>
            <span className="font-semibold text-slate-800">{selectedPreset.name}</span>
            <span className="text-slate-400 block text-[10px]">{selectedPreset.tag}</span>
          </div>
        </div>

        <Canvas shadows camera={{ position: [12, 10, 12], fov: 45 }}>
          <color attach="background" args={["#f8fafc"]} />
          <ambientLight intensity={0.65} />
          <directionalLight
            castShadow
            position={[15, 25, 10]}
            intensity={1.6}
            shadow-mapSize={[2048, 2048]}
            shadow-bias={-0.0001}
          />
          <directionalLight position={[-15, 12, -10]} intensity={0.4} />
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
                <GLBBuildingAsset
                  key={selectedPreset.id}
                  modelPath={selectedPreset.modelPath}
                  scale={selectedPreset.scale}
                  targetY={0}
                />
              </Suspense>

              {/* Glowing base ring and boundary */}
              <mesh position={[0, 0.03, 0]} rotation={[-Math.PI / 2, 0, 0]}>
                <ringGeometry args={[2.2, 2.5, 32]} />
                <meshBasicMaterial color="#f97316" side={THREE.DoubleSide} />
              </mesh>

              {/* 3D Floating Billboard Tag */}
              <Html position={[0, 4.2, 0]} center distanceFactor={14}>
                <div className="px-3 py-1.5 rounded-xl text-xs font-bold whitespace-nowrap shadow-xl border bg-slate-900/90 backdrop-blur text-white border-white/20 pointer-events-none flex items-center gap-2">
                  <span>{selectedPreset.previewIcon}</span>
                  <div>
                    <span>{currentBuildingName}</span>
                    <span className="block text-[10px] text-orange font-normal">
                      {currentFloors} Floors · {selectedPreset.category}
                    </span>
                  </div>
                </div>
              </Html>
            </group>
          )}

          <OrbitControls
            enableDamping
            dampingFactor={0.05}
            maxPolarAngle={Math.PI / 2 - 0.05}
            minDistance={4}
            maxDistance={40}
          />
          <ContactShadows position={[0, 0.02, 0]} opacity={0.45} scale={40} blur={2} far={10} />
        </Canvas>
      </div>

      {/* Docked Form Panel (Right Sidebar) */}
      <div className="w-full lg:w-96 bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex flex-col justify-between overflow-y-auto">
        <div className="space-y-5">
          <div className="flex items-center gap-2 pb-4 border-b border-slate-100">
            <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center text-orange">
              <Building2 size={18} />
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-base">Building Setup</h3>
              <p className="text-xs text-slate-500">Pick 3D model & define properties</p>
            </div>
          </div>

          {/* 3D Model Preset Selector */}
          <div className="space-y-2">
            <Label className="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
              <Sparkles size={13} className="text-orange" />
              Choose 3D Building Model
            </Label>
            <div className="grid grid-cols-1 gap-2.5">
              {BUILDING_PRESETS.map((preset) => {
                const isSelected = selectedPreset.id === preset.id;
                return (
                  <button
                    key={preset.id}
                    type="button"
                    onClick={() => handleSelectPreset(preset)}
                    className={`p-3 rounded-xl border text-left transition-all flex items-start justify-between gap-3 ${
                      isSelected
                        ? "border-orange bg-orange/5 shadow-sm ring-2 ring-orange/20"
                        : "border-slate-200 hover:border-slate-300 hover:bg-slate-50"
                    }`}
                  >
                    <div className="flex items-start gap-3">
                      <span className="text-2xl pt-0.5">{preset.previewIcon}</span>
                      <div>
                        <p className={`text-xs font-bold ${isSelected ? "text-orange" : "text-slate-800"}`}>
                          {preset.name}
                        </p>
                        <p className="text-[11px] text-slate-500">{preset.category}</p>
                        <span className="inline-block text-[10px] px-2 py-0.5 mt-1 rounded bg-slate-100 text-slate-600 font-mono">
                          {preset.tag}
                        </span>
                      </div>
                    </div>
                    {isSelected && (
                      <div className="w-5 h-5 rounded-full bg-orange text-white flex items-center justify-center shrink-0">
                        <Check size={12} strokeWidth={3} />
                      </div>
                    )}
                  </button>
                );
              })}
            </div>
          </div>

          <form id="building-creator-form" onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4 pt-2">
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
                placeholder="e.g. Tower Alpha, Grand Hotel"
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
            <div className="p-3 bg-slate-50 rounded-lg border border-slate-200 text-xs space-y-1 mt-3">
              <div className="flex justify-between">
                <span className="text-slate-500">Spatial Position:</span>
                <span className="font-mono font-medium text-slate-800">
                  [{position[0]}, 0, {position[2]}]
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">Active 3D Model:</span>
                <span className="font-semibold text-slate-800 flex items-center gap-1">
                  <span>{selectedPreset.previewIcon}</span> {selectedPreset.tag}
                </span>
              </div>
            </div>
          </form>
        </div>

        {/* Submit action button */}
        <div className="pt-5 border-t border-slate-100 mt-4">
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
            {isSubmitting ? "Saving Building..." : "Save Building"}
          </Button>
        </div>
      </div>
    </div>
  );
}
