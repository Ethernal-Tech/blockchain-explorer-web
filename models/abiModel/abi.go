package abiModel

type AbiItem struct {
	Anonymous bool `json:"anonymous"`
	Inputs    []struct {
		Indexed      bool   `json:"indexed"`
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	Name    string `json:"name"`
	Outputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs"`
	StateMutability string `json:"stateMutability"`
	Type            string `json:"type"`
}

type EventItem struct {
	Anonymous bool `json:"anonymous"`
	Inputs    []struct {
		Indexed      bool   `json:"indexed"`
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	Name string `json:"name"`
	Type string `json:"type"`
}
