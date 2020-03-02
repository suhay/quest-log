package types

// RequirementResolver resolves any requirement arrays
type RequirementResolver struct {
	Requirement *Requirement
}

// Requirement is something that must be resolved by the hook before it is considered finished
type Requirement struct {
	Key   string
	Value string
	Or    []Requirement
}

// Key is variable name to evaluate
func (r *RequirementResolver) Key() *string {
	return &r.Requirement.Key
}

// Value is the exact value the evaluated variable must equal
func (r *RequirementResolver) Value() *string {
	return &r.Requirement.Value
}

// Or returns a nested Requirements array where any within the array can qualify for completion
func (r *RequirementResolver) Or() *[]*RequirementResolver {
	var resolvers []*RequirementResolver

	for _, e := range r.Requirement.Or {
		resolver := RequirementResolver{&e}
		resolvers = append(resolvers, &resolver)
	}

	return &resolvers
}
