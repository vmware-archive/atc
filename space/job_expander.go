package space

import (
	"reflect"

	"github.com/concourse/atc"
)

type JobExpander struct {
	ResourceSpaces map[string][]string
}

type Permutation map[string]string

func FindCombinations(resourceSpaces atc.SpaceConfig) []Permutation {
	permutations := []Permutation{}

	for resource, spaces := range resourceSpaces {
		otherResourceSpaces := atc.SpaceConfig{}
		for otherResource, otherSpaces := range resourceSpaces {
			if otherResource == resource {
				continue
			}

			otherResourceSpaces[otherResource] = otherSpaces
		}

		otherCombinations := FindCombinations(otherResourceSpaces)

		for _, space := range spaces {
			withMe := Permutation{resource: space}

			if len(otherCombinations) == 0 {
				permutations = append(permutations, withMe)
			} else {
				for _, permutation := range otherCombinations {
					dup := Permutation{resource: space}
					for r, s := range permutation {
						dup[r] = s
					}

					var exists bool
					for _, existing := range permutations {
						if reflect.DeepEqual(existing, dup) {
							exists = true
							break
						}
					}

					if !exists {
						permutations = append(permutations, dup)
					}
				}
			}
		}
	}

	return permutations
}
