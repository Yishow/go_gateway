# Frontend-Backend Integration Map Template

For each screen, document:

- route
- page entry
- primary domain objects
- current integration chain
- whether the current contract is reusable unchanged
- whether an adapter is needed
- verification plan

## Integration chain format

`route -> page -> hook -> service -> API`

## Example

### Screen

- route: `/devices`
- page entry: `DevicesPage`
- primary domain objects: device, status, capability

### Current integration chain

- route -> `/devices`
- page -> `DevicesPage`
- hook -> `useDevicesQuery`
- service -> `deviceAPI`
- API -> `/api/devices`

### Reuse status

- partial

### Adapter

- needed: yes
- name: `DeviceSummaryAdapter`
- purpose: reshape backend payload into the redesigned screen summary model

### Verification

- typecheck
- build
- targeted tests
- route sanity check
