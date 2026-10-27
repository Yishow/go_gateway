package dbtarget

import "errors"

// classifySQLServerInsertError uses the driver's stable error number through
// wrappers, so neither server language nor rejected row contents affect retry.
func classifySQLServerInsertError(err error) (InsertErrorClass, bool) {
	var numbered interface{ SQLErrorNumber() int32 }
	if !errors.As(err, &numbered) {
		return "", false
	}
	switch numbered.SQLErrorNumber() {
	case 207, 208, 229:
		return InsertErrorTarget, true
	case 2601, 2627, 2628, 515, 547, 8115:
		return InsertErrorRow, true
	default:
		// Deadlock (1205), lock timeout (1222), transport failures and unknown
		// engine failures retain the payload for the existing bounded retry.
		return InsertErrorTransient, true
	}
}
