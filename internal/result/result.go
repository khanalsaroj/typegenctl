package result

type ValidationResult struct {
	Success []string `json:"success,omitempty"`
	Info    []string `json:"info,omitempty"`
	Failure []error  `json:"failure,omitempty"`
}

func (r *ValidationResult) AddSuccess(msg string) {
	r.Success = append(r.Success, msg)
}

func (r *ValidationResult) AddInfo(msg string) {
	r.Info = append(r.Info, msg)
}

func (r *ValidationResult) AddFailure(err error) {
	r.Failure = append(r.Failure, err)
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Failure) > 0
}
