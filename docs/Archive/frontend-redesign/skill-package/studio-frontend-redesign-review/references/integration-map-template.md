# Frontend-Backend Integration Map Template

## Goal

Map current frontend screens to backend contracts and define how AI-generated frontend output reconnects safely.

## Mapping format

### Screen

- route:
- page entry:
- primary domain objects:

### Current integration chain

- page ->
- hook ->
- service ->
- API ->

### Can reuse unchanged?

- yes / partial / no

### Adapter needed?

- yes / no
- adapter name:
- purpose:

### Verification

- typecheck:
- build:
- tests:
- manual route sanity:

## Example

### Screen

- route: `/studio`
- page entry: `DatalinkWorkbenchPage`
- primary domain objects: device, source rule, point, tag, mapping

### Current integration chain

- page -> workbench subpages
- hook -> `useDevicesQuery`
- service -> `deviceAPI`
- API -> `/api/v1/datalink/devices`

### Can reuse unchanged?

- partial

### Adapter needed?

- yes
- adapter name: `DeviceCapabilityAdapter`
- purpose: reshape backend device payload into the new workflow shell summary

### Verification

- typecheck: required
- build: required
- tests: targeted device/workbench tests
- manual route sanity: `/studio` device flow
