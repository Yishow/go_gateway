package recordingplan

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func recreateScope(policy *bool) SchemaPreviewScope {
	scope := groupPreviewScope()
	scope.RecreateMissing = policy
	return scope
}

func flag(value bool) *bool { return &value }

// TestSchemaPreviewShowsAndFreezesRecreateMissing: the token a user confirms
// shows the table-recovery policy, the policy is part of what the digest
// protects, and a policy that changed since the preview makes the confirmation
// stale. A layout without a policy keeps the digest it always had.
func TestSchemaPreviewShowsAndFreezesRecreateMissing(t *testing.T) {
	prepare := func(t *testing.T, policy *bool) *SchemaPreviewToken {
		t.Helper()
		svc, _ := newPreviewService()
		token, err := svc.PrepareSchemaPreview(t.Context(), recreateScope(policy), (&inspectorStub{}).inspect)
		require.NoError(t, err)
		return token
	}

	on, off, unset := prepare(t, flag(true)), prepare(t, flag(false)), prepare(t, nil)
	require.Equal(t, flag(true), on.RecreateMissing, "the preview shows the policy it was prepared with")
	require.Equal(t, flag(false), off.RecreateMissing)
	require.Nil(t, unset.RecreateMissing, "a layout without a policy shows none")

	digestOf := func(token *SchemaPreviewToken) string {
		digest, err := previewDigest(token)
		require.NoError(t, err)
		return digest
	}
	require.Equal(t, on.Digest, digestOf(on), "the digest covers the stored token")
	forged := *on
	forged.RecreateMissing = flag(false)
	require.NotEqual(t, on.Digest, digestOf(&forged), "flipping the policy of a stored token changes its digest")
	forged.RecreateMissing = nil
	require.NotEqual(t, on.Digest, digestOf(&forged), "so does removing it")

	legacy := *unset
	legacy.OperationID, legacy.Token = on.OperationID, on.Token
	withoutField := legacy
	withoutField.RecreateMissing = nil
	require.Equal(t, digestOf(&legacy), digestOf(&withoutField), "an absent policy adds nothing to the digest")

	t.Run("a policy that changed since the preview is stale", func(t *testing.T) {
		for name, current := range map[string]*bool{"turned off": flag(false), "removed": nil} {
			svc, _ := newPreviewService()
			token, err := svc.PrepareSchemaPreview(t.Context(), recreateScope(flag(true)), (&inspectorStub{}).inspect)
			require.NoError(t, err)
			_, err = svc.ValidatePreviewTokenForApply(t.Context(), token.Token, recreateScope(current))
			require.True(t, errors.Is(err, ErrPreviewTokenStale), "%s: %v", name, err)
		}
		svc, _ := newPreviewService()
		token, err := svc.PrepareSchemaPreview(t.Context(), recreateScope(flag(true)), (&inspectorStub{}).inspect)
		require.NoError(t, err)
		_, err = svc.ValidatePreviewTokenForApply(t.Context(), token.Token, recreateScope(flag(true)))
		require.NoError(t, err, "the same policy still confirms")
	})
}

// TestStoredPreviewKeepsRecreateMissing: the policy a preview showed survives
// the trip through the token store, so the digest the user confirmed still
// verifies when the confirmation arrives.
func TestStoredPreviewKeepsRecreateMissing(t *testing.T) {
	for name, policy := range map[string]*bool{"on": flag(true), "off": flag(false), "absent": nil} {
		t.Run(name, func(t *testing.T) {
			svc := NewService(NewSQLRepository(recordingRepositoryDB(t)))
			token, err := svc.PrepareSchemaPreview(t.Context(), recreateScope(policy), (&inspectorStub{}).inspect)
			require.NoError(t, err)
			stored, err := svc.PreviewTokenForWorkspace(t.Context(), token.WorkspaceID, token.Token)
			require.NoError(t, err)
			require.Equal(t, policy, stored.RecreateMissing)
			digest, err := previewDigest(stored)
			require.NoError(t, err)
			require.Equal(t, token.Digest, digest, "the stored token still carries the digest that was shown")
		})
	}
}
