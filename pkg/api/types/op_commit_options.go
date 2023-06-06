package types

type CommitOperationOptions struct {
	Req    string
	Rawref string

	Author  string
	Message string
	Change  []string
}
