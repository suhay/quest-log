package types

// HookResolver resolves the Hook type
type HookResolver struct {
	H *Hook
}

// Hook are the main parts of a Thread of Entry
type Hook struct {
	Hook        string
	Trigger     string
	Closing     string
	Required    []Requirement
	Tags        []*string
	Event       string
}

// Hook returns the hook text
func (r *HookResolver) Hook() *string {
	return &r.H.Hook
}

// Trigger returns the expression that can trigger this particular hook to fire
func (r *HookResolver) Trigger() *string {
	return &r.H.Trigger
}

// Closing returns the expression that is evaluated once the hook's goal is achieved
func (r *HookResolver) Closing() *string {
	return &r.H.Closing
}

// Event returns the expression that is evaluated once the hook's begins
func (r *HookResolver) Event() *string {
	return &r.H.Event
}

// Required returns an array of key values that must be achieved in order to complete the hook
func (r *HookResolver) Required() *[]*RequirementResolver {
	var resolvers []*RequirementResolver

	for _, e := range r.H.Required {
		resolver := RequirementResolver{&e}
		resolvers = append(resolvers, &resolver)
	}

	return &resolvers
}

// Tags are used to relate hooks together
func (r *HookResolver) Tags() *[]*string {
	return &r.H.Tags
}
