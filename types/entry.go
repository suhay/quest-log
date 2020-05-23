package types

// EntryResolver resolves any Entry objects
type EntryResolver struct {
	Entry   *Entry
	Entries []Entry
}

// Entry is a collection of related thread hooks
type Entry struct {
	Thread
	Closing string
}

// Name is the database name of this particular entry
func (r *EntryResolver) Name() string {
	return r.Entry.Name
}

// Hooks returns the quest hooks
func (r *EntryResolver) Hooks() *[]*HookResolver {
	var resolvers []*HookResolver

	for _, e := range r.Entry.Hooks {
		resolver := HookResolver{&e}
		resolvers = append(resolvers, &resolver)
	}

	return &resolvers
}

// Closing is the expression to be evaluated after the Goal has been completed
func (r *EntryResolver) Closing() *string {
	return &r.Entry.Closing
}

// Tags relate entries together and also are used for find child threads to include
func (r *EntryResolver) Tags() *[]*string {
	return &r.Entry.Tags
}

// Perspective is the written perspective the entry should contain threads written in
func (r *EntryResolver) Perspective() *string {
	return &r.Entry.Perspective
}
