package algorithm

import (
	"fmt"
	"sort"
	"strings"
)

type JobPermutationSet map[int]struct{}

func (set JobPermutationSet) Add(id int) {
	set[id] = struct{}{}
}

func (set JobPermutationSet) Contains(jobID int) bool {
	_, found := set[jobID]
	return found
}

func (set JobPermutationSet) Union(otherSet JobPermutationSet) JobPermutationSet {
	newSet := JobPermutationSet{}

	for jobID, _ := range set {
		newSet[jobID] = struct{}{}
	}

	for jobID, _ := range otherSet {
		newSet[jobID] = struct{}{}
	}

	return newSet
}

func (set JobPermutationSet) Intersect(otherSet JobPermutationSet) JobPermutationSet {
	result := JobPermutationSet{}

	for key, val := range set {
		_, found := otherSet[key]
		if found {
			result[key] = val
		}
	}

	return result
}

func (set JobPermutationSet) Equal(otherSet JobPermutationSet) bool {
	if len(set) != len(otherSet) {
		return false
	}

	for x, _ := range set {
		if !otherSet.Contains(x) {
			return false
		}
	}

	return true
}

func (set JobPermutationSet) String() string {
	xs := []string{}
	for x, _ := range set {
		xs = append(xs, fmt.Sprintf("%v", x))
	}

	sort.Strings(xs)

	return fmt.Sprintf("{%s}", strings.Join(xs, " "))
}
