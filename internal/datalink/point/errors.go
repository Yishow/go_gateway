package point

import "errors"

// ErrPointNotFound 表示依 ID 查無點位（含已刪除）。
var ErrPointNotFound = errors.New("point not found")
