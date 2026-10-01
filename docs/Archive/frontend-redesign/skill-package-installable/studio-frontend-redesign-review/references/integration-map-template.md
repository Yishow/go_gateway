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

- route: `/studio`
- page entry: `DatalinkWorkbenchPage`
- primary domain objects: device, source rule, point, tag, mapping

### Current integration chain

- route -> `/studio`
- page -> `DatalinkWorkbenchPage`
- hook -> `useDevicesQuery`
- service -> `deviceAPI`
- API -> `/api/v1/datalink/devices`

### Reuse status

- partial

### Adapter

- needed: yes
- name: `DeviceCapabilityAdapter`
- purpose: reshape backend device payload into the new workflow shell summary

### Verification

- typecheck
- build
- targeted tests
- route sanity check
