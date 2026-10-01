# API Inventory Template

For each API group, document:

- endpoint
- method
- screen consumers
- current frontend service wrapper
- current hook consumers
- redesign impact

## Recommended groups

- `/api/v1/test/*`
- `/api/v1/debug/*`
- `/api/v1/datalink/*`

## Example

### Group: `/api/v1/datalink/devices`

- Endpoints:
  - `GET /devices`
  - `POST /devices`
  - `GET /devices/:id`
  - `PUT /devices/:id`
  - `DELETE /devices/:id`
  - `POST /devices/:id/test`
- Used by:
  - `/studio` Device step
- Frontend wrapper:
  - `deviceAPI`
- Hooks:
  - `useDevicesQuery`
  - `useCreateDeviceMutation`
  - `useUpdateDeviceMutation`
- Redesign impact:
  - likely reusable unchanged
