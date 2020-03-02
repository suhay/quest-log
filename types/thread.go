package types

// ThreadResolver resolves any Thread objects
type ThreadResolver struct {
	Thread *Thread
}

// Thread is a collection of related thread hooks
type Thread struct {
	Name        string
	Hooks       []Hook
	Tags        []*string
	Perspective string
}

// Name is the database name of this particular thread
func (r *ThreadResolver) Name() string {
	return r.Thread.Name
}

// Hooks returns the quest hooks
func (r *ThreadResolver) Hooks() *[]*HookResolver {
	var resolvers []*HookResolver

	for _, e := range r.Thread.Hooks {
		resolver := HookResolver{&e}
		resolvers = append(resolvers, &resolver)
	}

	return &resolvers
}

// Tags are used to relate one thread to another
func (r *ThreadResolver) Tags() *[]*string {
	return &r.Thread.Tags
}

// Perspective is the perspective the thread should pull hooks in as
func (r *ThreadResolver) Perspective() *string {
	return &r.Thread.Perspective
}
