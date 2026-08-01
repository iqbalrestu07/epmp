# EPMP SDK Schemas

This directory contains the unified YAML configuration files used by both the **Backend** and **Frontend** codegen tools.

By keeping the configuration unified, we ensure that both the backend database/REST APIs and the frontend UI/API clients remain in perfect sync for any given module.

## Usage

```bash
# Generate Backend Module
cd ../be/codegen
./epmp-codegen --config ../../schemas/property.yaml --output ../../../../backend/internal

# Generate Frontend Module
cd ../fe/codegen
./epmp-fe-codegen --config ../../schemas/property.yaml --output ../../../../frontend/src
```
