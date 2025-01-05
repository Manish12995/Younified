package graphqlclient

import (
	"fmt"
	"strings"
)

// MutationBuilder helps construct GraphQL mutations
type MutationBuilder struct {
	name      string
	fields    []string
	input     map[string]interface{}
	fragments []string
	inputType string
}

// NewMutationBuilder creates a new mutation builder
func NewMutationBuilder() *MutationBuilder {
	return &MutationBuilder{
		name:      "",
		input:     make(map[string]interface{}),
		fields:    []string{},
		inputType: "",
	}
}

// SetMutationName allows dynamically changing the mutation name after initial creation
func (mb *MutationBuilder) SetMutationName(name string) *MutationBuilder {
	mb.name = name
	return mb
}

// SetInputName allows dynamically changing the input type name
func (mb *MutationBuilder) SetInputName(inputType string) *MutationBuilder {
	mb.inputType = inputType
	return mb
}

// SetInput sets the input for the mutation
func (mb *MutationBuilder) SetInput(input map[string]interface{}) *MutationBuilder {
	mb.input = input
	return mb
}

// AddField adds a return field to the mutation
func (mb *MutationBuilder) AddField(field string) *MutationBuilder {
	mb.fields = append(mb.fields, field)
	return mb
}

// AddFragment adds a fragment to the mutation
func (mb *MutationBuilder) AddFragment(fragment string) *MutationBuilder {
	mb.fragments = append(mb.fragments, fragment)
	return mb
}

// Build constructs the final GraphQL mutation
func (mb *MutationBuilder) Build() (string, map[string]interface{}) {
	variables := map[string]interface{}{
		"input": make(map[string]interface{}),
	}

	for k, v := range mb.input {
		variables["input"].(map[string]interface{})[k] = v
	}
	var mutation string
	//check if the mutation has fields {} in return
	if len(mb.fields) > 0 {
		mutation = fmt.Sprintf(`
            mutation %s($input: %s!) {
                %s(input: $input) {
                    %s
                }
            }
        `,
			mb.name,
			mb.inputType,
			mb.name,
			strings.Join(mb.fields, " "),
		)
	} else {
		mutation = fmt.Sprintf(`
            mutation %s($input: %s!) {
                %s(input: $input)
            }
        `,
			mb.name,
			mb.inputType,
			mb.name,
		)
	}

	return mutation, variables
}
