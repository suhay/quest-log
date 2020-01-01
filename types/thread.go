package types

// ThreadResolver resolves any Thread objects
type ThreadResolver struct {
	E *Thread
}

// Thread is a collection of related thread hooks
type Thread struct {
	Name        string
	Plot        *Hook
	Climax      *Hook
	Goal        *Hook
	Tags        []*string
	Perspective string
}

// Name is the database name of this particular thread
func (r *ThreadResolver) Name() string {
	return r.E.Name
}

// Plot returns the plot hook, usually the first hook in a thread
func (r *ThreadResolver) Plot() *HookResolver {
	return &HookResolver{
		H: r.E.Plot,
	}
}

// Climax returns the climax hook, usually the middle of a thread which is activated once the Plot is considered completed
func (r *ThreadResolver) Climax() *HookResolver {
	return &HookResolver{
		H: r.E.Climax,
	}
}

// Goal returns the goal hook, usually the last part of the thread which is activated once the Climax is considered completed, or the Plot is if there is not climax
func (r *ThreadResolver) Goal() *HookResolver {
	return &HookResolver{
		H: r.E.Goal,
	}
}

// Tags are used to relate one thread to another
func (r *ThreadResolver) Tags() *[]*string {
	return &r.E.Tags
}

// Perspective is the perspective the thread should pull hooks in as
func (r *ThreadResolver) Perspective() *string {
	return &r.E.Perspective
}