package main

type Scope string
type Scopes []string

// IsAuthorized takes a slice of scopes (i.e. that an API token has) and
// checks if they sufficiently cover the required scopes
func (required Scopes) IsAuthorized(have Scopes) bool {

	difference := required.subtract(have)

	// i.e. if "have" contains every element of "required"
	// then required - have == []
	return len(difference) == 0
}

// [A, B, C] - [B, C] == [A]
func (scopes Scopes) subtract(scopesToRemove Scopes) Scopes {

	difference := scopes

	for _, item := range scopesToRemove {
		difference = difference.remove(Scope(item))
	}

	return difference
}

// [A, B, C] - [B] == [A, C]
func (scopes Scopes) remove(scope Scope) Scopes {
	var newScopes Scopes

	for _, i := range scopes {
		if Scope(i) != scope {
			newScopes = append(newScopes, i)
		}
	}

	return newScopes
}
