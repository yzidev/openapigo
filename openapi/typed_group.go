package openapi

// TypedGroup is a small helper to apply shared HandlerOption(s) (like WithTags/WithSecurity)
// once for multiple typed handlers (GETT/POSTT/etc).
//
// It is adapter-agnostic: adapters can wrap it with their own typed handler signatures.
type TypedGroup struct {
	opts []HandlerOption
}

func NewTypedGroup(opts ...HandlerOption) *TypedGroup {
	return &TypedGroup{opts: opts}
}

func (g *TypedGroup) Options() []HandlerOption {
	out := make([]HandlerOption, len(g.opts))
	copy(out, g.opts)
	return out
}

func (g *TypedGroup) With(opts ...HandlerOption) *TypedGroup {
	out := make([]HandlerOption, 0, len(g.opts)+len(opts))
	out = append(out, g.opts...)
	out = append(out, opts...)
	return &TypedGroup{opts: out}
}
