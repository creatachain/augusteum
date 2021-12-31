package client

// MSMQueryOptions can be used to provide options for MSMQuery call other
// than the DefaultMSMQueryOptions.
type MSMQueryOptions struct {
	Height int64
	Prove  bool
}

// DefaultMSMQueryOptions are latest height (0) and prove false.
var DefaultMSMQueryOptions = MSMQueryOptions{Height: 0, Prove: false}
