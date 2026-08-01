#!/bin/bash
set -e
cd /Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/be/codegen
./epmp-codegen --config ../../schemas/building.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/zone.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/bed.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/facility.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/roomtype.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/tenantidentity.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/tenantcontact.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/tenantdocument.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/asset.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/assetassignment.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/assetinspection.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/workorder.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/technician.yaml --output ../../../../backend/internal
./epmp-codegen --config ../../schemas/vendor.yaml --output ../../../../backend/internal

cd /Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/fe/codegen
./epmp-fe-codegen --config ../../schemas/building.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/zone.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/bed.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/facility.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/roomtype.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/tenantidentity.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/tenantcontact.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/tenantdocument.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/asset.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/assetassignment.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/assetinspection.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/workorder.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/technician.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/vendor.yaml --output ../../../../frontend/src
