package sourcerule

// MappingCandidateID returns the canonical identity of the mapping candidate
// derived from a rule/address pair. It is an identity, not activation proof.
func MappingCandidateID(ruleID, address string) (string, error) {
	return candidateID(buildMappingCandidateIdentity(ruleID, address, ""))
}
