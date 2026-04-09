# API Inventory Template

For each API group, document:

- endpoint
- method
- screen consumers
- current frontend service wrapper
- current hook or data-layer consumers
- redesign impact

## Recommended groups

- product APIs
- utility or debug APIs
- auth or session APIs
- upload or export APIs
- stream or realtime APIs

## Example

### Group: `/api/devices`

- Endpoints:
  - `GET /devices`
  - `POST /devices`
  - `GET /devices/:id`
  - `PUT /devices/:id`
  - `DELETE /devices/:id`
- Used by:
  - product device screen
- Frontend wrapper:
  - `deviceAPI`
- Hooks:
  - `useDevicesQuery`
  - `useCreateDeviceMutation`
- Redesign impact:
  - likely reusable unchanged
