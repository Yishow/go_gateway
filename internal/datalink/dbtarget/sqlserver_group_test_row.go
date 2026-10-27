package dbtarget

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	sqlServerGroupLongTablePattern = regexp.MustCompile(`^gw_long_[0-9a-f]{16}$`)
	sqlServerGroupLongOwnerPattern = regexp.MustCompile(`^_gw_owner_[0-9a-f]{24}$`)
)

func validateSQLServerGroupLongTestRowRef(ref GroupLongTestRowRef) error {
	invalid := func(reason string) error {
		return fmt.Errorf("%w: SQL Server long test %s", ErrOwnedRowRefInvalid, reason)
	}
	if !sqlServerGroupLongTablePattern.MatchString(ref.TableName) || !sqlServerGroupLongOwnerPattern.MatchString(ref.OwnerColumn) ||
		ref.ProvenanceColumn != longColProvenance || !sqlServerLongTestIdentity(ref.SchemaName, 128) {
		return invalid("canonical table scope")
	}
	for _, value := range []string{ref.WorkspaceID, ref.GroupID, ref.GroupRevision, ref.OperationOwner} {
		if !sqlServerLongTestIdentity(value, 191) {
			return invalid("identity is not exact Unicode text")
		}
	}
	for index, sampleID := range ref.SampleIDs {
		if !sqlServerLongTestIdentity(sampleID, 191) || sampleID != fmt.Sprintf("%s-%d", ref.OperationOwner, index) {
			return invalid("sample identity is not exact Unicode text")
		}
		if !groupLongRecordIDPattern.MatchString(ref.EffectKeys[index]) ||
			!groupLongRecordIDPattern.MatchString(ref.EffectDigests[ref.EffectKeys[index]]) {
			return invalid("receipt identity is not a canonical digest")
		}
	}
	return nil
}

func sqlServerLongTestIdentity(value string, maxUTF16 int) bool {
	return value != "" && utf8.ValidString(value) && !strings.HasSuffix(value, " ") &&
		!strings.ContainsRune(value, '\x00') && len(utf16.Encode([]rune(value))) <= maxUTF16
}
