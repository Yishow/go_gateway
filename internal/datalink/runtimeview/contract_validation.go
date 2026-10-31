package runtimeview

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var safeID = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$`)
var safeLabel = regexp.MustCompile(`^[\p{L}\p{N}_ -]{1,128}$`)
var utcTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9})?Z$`)

// ParseRuleEvidence is pure, fail-closed structural validation, NOT API authorization or proof
// verification. A future trusted adapter must establish all claimed evidence from its authorities.
func ParseRuleEvidence(data []byte, scope Scope) (*RuleEvidenceView, error) {
	if err := checkJSON(data, reflect.TypeFor[RuleEvidenceView](), ""); err != nil {
		return nil, errUnavailable
	}
	var v RuleEvidenceView
	if json.Unmarshal(data, &v) != nil {
		return nil, errUnavailable
	}
	if err := validateView(&v, scope); err != nil {
		return nil, err
	}
	return &v, nil
}

// checkJSON uses only declared field types/tags. RawMessage exists solely at the JSON boundary,
// never in a DTO. Duplicate, unknown, wrong-case keys, null primitives and missing fields fail.
func checkJSON(raw []byte, typ reflect.Type, tag reflect.StructTag) error {
	absent := len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
	if typ.Kind() == reflect.Pointer {
		if absent {
			return nil
		}
		return checkJSON(raw, typ.Elem(), tag)
	}
	if absent {
		return errUnavailable
	}
	switch typ.Kind() {
	case reflect.Struct:
		d := json.NewDecoder(bytes.NewReader(raw))
		tok, err := d.Token()
		if err != nil || tok != json.Delim('{') {
			return errUnavailable
		}
		fields := make(map[string]reflect.StructField, typ.NumField())
		seen := make(map[string]bool)
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			fields[strings.Split(f.Tag.Get("json"), ",")[0]] = f
		}
		for d.More() {
			key, err := d.Token()
			if err != nil {
				return errUnavailable
			}
			name, ok := key.(string)
			if !ok {
				return errUnavailable
			}
			f, ok := fields[name]
			if !ok || seen[name] {
				return errUnavailable
			}
			seen[name] = true
			var value json.RawMessage
			if d.Decode(&value) != nil || checkJSON(value, f.Type, f.Tag) != nil {
				return errUnavailable
			}
		}
		if _, err = d.Token(); err != nil {
			return errUnavailable
		}
		if _, err = d.Token(); err != io.EOF {
			return errUnavailable
		}
		for name, f := range fields {
			if !seen[name] && f.Type.Kind() != reflect.Pointer {
				return errUnavailable
			}
		}
	case reflect.Slice:
		var entries []json.RawMessage
		if json.Unmarshal(raw, &entries) != nil || entries == nil {
			return errUnavailable
		}
		for _, entry := range entries {
			if checkJSON(entry, typ.Elem(), tag) != nil {
				return errUnavailable
			}
		}
	default:
		value := reflect.New(typ)
		if json.Unmarshal(raw, value.Interface()) != nil {
			return errUnavailable
		}
		v := value.Elem()
		switch typ.Kind() {
		case reflect.String:
			s := v.String()
			enums := tag.Get("enum")
			if typ == reflect.TypeFor[Reason]() {
				enums = reasonValues
			}
			if enums != "" {
				if !member(s, enums) {
					return errUnavailable
				}
				return nil
			}
			switch tag.Get("format") {
			case "time":
				if !utcTime.MatchString(s) {
					return errUnavailable
				}
				if _, err := time.Parse(time.RFC3339Nano, s); err != nil {
					return errUnavailable
				}
			case "label":
				if !safeLabel.MatchString(s) {
					return errUnavailable
				}
			default:
				if !safeID.MatchString(s) {
					return errUnavailable
				}
			}
		case reflect.Int, reflect.Int64, reflect.Float64:
			n := float64(0)
			if typ.Kind() == reflect.Float64 {
				n = v.Float()
			} else {
				n = float64(v.Int())
			}
			if math.IsNaN(n) || math.IsInf(n, 0) || math.Abs(n) > 9007199254740991 {
				return errUnavailable
			}
			min := float64(0)
			if tag.Get("signed") == "true" {
				min = -9007199254740991
			}
			if s := tag.Get("min"); s != "" {
				min, _ = strconv.ParseFloat(s, 64)
			}
			max := float64(9007199254740991)
			if s := tag.Get("max"); s != "" {
				max, _ = strconv.ParseFloat(s, 64)
			}
			if n < min || n > max {
				return errUnavailable
			}
		case reflect.Bool:
		default:
			return errUnavailable
		}
	}
	return nil
}

func checkBinding(b Binding, revision *string, scope Scope) error {
	if b.WorkspaceID != scope.WorkspaceID {
		return errForeign
	}
	if b.WorkspaceRevision != scope.WorkspaceRevision || revision == nil || b.EntityRevision != *revision {
		return errStale
	}
	return nil
}
func checkTransform(s TransformStep) bool {
	switch s.Kind {
	case "scale":
		return s.Multiplier != nil && s.Offset != nil && s.Order == nil && s.TargetType == nil
	case "cast":
		return s.TargetType != nil && s.Multiplier == nil && s.Offset == nil && s.Order == nil
	case "byte-order":
		return s.Order != nil && s.Multiplier == nil && s.Offset == nil && s.TargetType == nil
	default:
		return s.Order == nil && s.Multiplier == nil && s.Offset == nil && s.TargetType == nil
	}
}
func validateSource(s *SourceView, scope Scope) error {
	if s == nil {
		return nil
	}
	for _, pair := range []struct {
		facts    *SourceFacts
		revision *string
	}{{s.Saved, s.SavedRevision}, {s.Effective, s.EffectiveRevision}} {
		if pair.facts == nil {
			continue
		}
		if err := checkBinding(pair.facts.Binding, pair.revision, scope); err != nil {
			return err
		}
		for _, step := range pair.facts.Transforms {
			if !checkTransform(step) {
				return errUnavailable
			}
		}
	}
	if s.State == "saved-only" && (s.Saved == nil || s.Effective != nil) {
		return errUnavailable
	}
	if s.State != "effective" && s.Effective != nil {
		return errUnavailable
	}
	unsupported := false
	if s.Effective != nil {
		for _, step := range s.Effective.Transforms {
			unsupported = unsupported || step.Kind == "unsupported"
		}
	}
	if s.State == "effective" && (unsupported || s.Proof != "immutable-source-facts" || s.Effective == nil || s.EffectiveRevision == nil || s.Effective.MappingRevision == nil || s.Effective.ReadType == nil || s.Effective.ResultType == nil || s.Effective.Address == nil || s.Effective.Polling == nil) {
		s.State = "unavailable"
		s.Effective = nil
		s.Proof = "unavailable"
	}
	return nil
}

func validateView(v *RuleEvidenceView, scope Scope) error {
	if !safeID.MatchString(scope.WorkspaceID) || !safeID.MatchString(scope.WorkspaceRevision) || !safeID.MatchString(scope.DeviceID) {
		return errUnavailable
	}
	if v.WorkspaceID != scope.WorkspaceID || v.DeviceID != scope.DeviceID {
		return errForeign
	}
	if v.WorkspaceRevision != scope.WorkspaceRevision {
		return errStale
	}
	if v.Freshness != "unavailable" && v.ObservedAt == nil {
		return errUnavailable
	}
	if err := validateSource(v.Source, scope); err != nil {
		return err
	}
	groups := map[string]bool{}
	for i := range v.Groups {
		g := &v.Groups[i]
		if groups[g.GroupID] {
			return errUnavailable
		}
		groups[g.GroupID] = true
		if err := validateGroup(g, scope, v.ObservedAt); err != nil {
			return err
		}
	}
	if v.Locator != nil {
		if err := ValidateEntityLocator(*v.Locator, scope); err != nil {
			return err
		}
		l := v.Locator
		if (l.GroupID != nil && !groups[*l.GroupID]) || (l.RuleID != nil && (v.Source == nil || *l.RuleID != v.Source.RuleID)) || (l.MappingID != nil && (v.Source == nil || v.Source.MappingID == nil || *l.MappingID != *v.Source.MappingID)) {
			return errStale
		}
	}
	return nil
}

// ValidateEntityLocator checks structure and selected scope only. Deleted/ownership checks remain
// the future resolver's responsibility. Historical group_revision never authorizes a mutation.
func ValidateEntityLocator(l EntityLocator, scope Scope) error {
	raw, err := json.Marshal(l)
	if err != nil || checkJSON(raw, reflect.TypeFor[EntityLocator](), "") != nil {
		return errUnavailable
	}
	if l.WorkspaceID != scope.WorkspaceID || (l.DeviceID != nil && *l.DeviceID != scope.DeviceID) {
		return errForeign
	}
	switch l.Step {
	case 1:
		if l.DeviceID == nil || l.RuleID != nil || l.MappingID != nil || l.GroupID != nil || l.GroupRevision != nil {
			return errUnavailable
		}
	case 2:
		if l.DeviceID == nil || l.RuleID == nil || l.MappingID != nil || l.GroupID != nil || l.GroupRevision != nil {
			return errUnavailable
		}
	case 3:
		if l.DeviceID == nil || l.MappingID == nil || l.RuleID != nil || l.GroupID != nil || l.GroupRevision != nil {
			return errUnavailable
		}
	case 4:
		if l.GroupID == nil || l.RuleID != nil || l.MappingID != nil {
			return errUnavailable
		}
	}
	return nil
}
