package types

// EntryResolver resolves any Entry objects
type EntryResolver struct {
	E *Entry
}

// Entry is a collection of related thread hooks
type Entry struct {
	Name        string
	Plot        *Hook
	Climax      *Hook
	Goal        *Hook
	Tags        []*string
	Closing     string
	Perspective string
}

// Name is the database name of this particular entry
func (r *EntryResolver) Name() string {
	return r.E.Name
}

// Plot returns the plot hook, usually the first hook in a thread
func (r *EntryResolver) Plot() *HookResolver {
	return &HookResolver{
		H: r.E.Plot,
	}
}

// Climax returns the climax hook, usually the middle of a thread which is activated once the Plot is considered completed
func (r *EntryResolver) Climax() *HookResolver {
	return &HookResolver{
		H: r.E.Climax,
	}
}

// Goal returns the goal hook, usually the last part of the thread which is activated once the Climax is considered completed, or the Plot is if there is not climax
func (r *EntryResolver) Goal() *HookResolver {
	return &HookResolver{
		H: r.E.Goal,
	}
}

// Closing is the expression to be evaluated after the Goal has been completed
func (r *EntryResolver) Closing() *string {
	return &r.E.Closing
}

// Tags relate enteries together and also are used for find chil threads to include
func (r *EntryResolver) Tags() *[]*string {
	return &r.E.Tags
}

// Perspective is the written perspective the entry should contain threads written in
func (r *EntryResolver) Perspective() *string {
	return &r.E.Perspective
}