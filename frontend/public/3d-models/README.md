# 3D Model Assets

Place your 3D model files here. Supported formats: `.glb`, `.gltf`, `.obj`, `.fbx`.

## Directory Structure

```
public/3d-models/
├── buildings/    # 3D models for property/building exteriors
├── floors/       # 3D models for floor plans
├── rooms/        # 3D models for individual rooms
└── README.md     # This file
```

## Naming Convention

- **Buildings**: `{property_type}.glb` — e.g. `boarding_house.glb`, `apartment.glb`
- **Floors**: `floor_template.glb`
- **Rooms**: `room_{type}.glb` — e.g. `room_standard.glb`, `room_penthouse.glb`

## Usage

These models are loaded by the interactive 3D view using Three.js.
The app will look for models in `/3d-models/` at runtime.

If no model is found, a procedural placeholder (simple box geometry) will be rendered.
