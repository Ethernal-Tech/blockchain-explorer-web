package abiModel

type AbiItem struct {
	Anonymous bool `json:"anonymous"`
	Inputs    []struct {
		Components []struct {
			InternalType string `json:"internalType"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"components,omitempty"`
		Indexed      bool   `json:"indexed,omitempty"`
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
		Components []struct {
			InternalType string `json:"internalType"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"components,omitempty"`
		Indexed      bool   `json:"indexed,omitempty"`
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type FunctionItem struct {
	Inputs []struct {
		Components []struct {
			InternalType string `json:"internalType"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"components,omitempty"`
		Indexed      bool   `json:"indexed,omitempty"`
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	Outputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs"`
	Name            string `json:"name"`
	StateMutability string `json:"stateMutability"`
	Type            string `json:"type"`
}
