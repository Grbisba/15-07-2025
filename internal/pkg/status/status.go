package status

type Status string

const (
	Created    Status = "created"
	InProcess  Status = "in_process"
	Completed  Status = "completed"
	Downloaded Status = "downloaded"
	Unknown    Status = "unknown"

	Uploaded    Status = "uploaded"
	Failed      Status = "failed"
	BadURL      Status = "bad_url"
	Unsupported Status = "unsupported"
)
