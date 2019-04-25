package chainpaxos

const (
	OK             = "OK"
	ErrNoKey       = "ErrNoKey"
	ErrEnCodeError = "ErrEncodeError"
	LastBlock = -1
)

type Err string

// Put or Append
type PutAppendArgs struct {
	// You'll have to add definitions here.
	Position   int
	Value string
	Op    string // "Put" or "Append"
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.

	ClientId  int
	ClientSeq int
}

type PutAppendReply struct {
	Err Err
}

type GetArgs struct {
	Position int
	// You'll have to add definitions here.
	ClientId  int
	ClientSeq int
}

type GetReply struct {
	Err   Err
	Value string
}
