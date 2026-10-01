# API Inventory Template

## Backend route groups

- `/api/v1/test/*`
- `/api/v1/debug/*`
- `/api/v1/datalink/*`

## For each group, capture

- endpoint
- method
- primary screen consumers
- current frontend service wrapper
- current hook consumers
- redesign impact

## Recommended output format

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
- Notes:
  - stable contract
  - likely reusable unchanged
