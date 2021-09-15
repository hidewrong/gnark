// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiled

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/consensys/gnark/backend/hint"
)

// LinearExpression represent a linear expression of variables
type LinearExpression []Term

// Clone returns a copy of the underlying slice
func (l LinearExpression) Clone() LinearExpression {
	res := make(LinearExpression, len(l))
	copy(res, l)
	return res
}

// Len return the lenght of the LinearExpression (implements Sort interface)
func (l LinearExpression) Len() int {
	return len(l)
}

// Swap swaps terms in the LinearExpression (implements Sort interface)
func (l LinearExpression) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less returns true if variableID for term at i is less than variableID for term at j (implements Sort interface)
func (l LinearExpression) Less(i, j int) bool {
	_, iID, iVis := l[i].Unpack()
	_, jID, jVis := l[j].Unpack()
	if iVis == jVis {
		return iID < jID
	}
	return iVis > jVis
}

// R1C used to compute the wires
type R1C struct {
	L LinearExpression
	R LinearExpression
	O LinearExpression
}

// LogEntry is used as a shared data structure between the frontend and the backend
// to represent string values (in logs or debug info) where a value is not known at compile time
// (which is the case for variables that need to be resolved in the R1CS)
type LogEntry struct {
	Format    string
	ToResolve []int
}

// Visibility encodes a Variable (or wire) visibility
// Possible values are Unset, Internal, Secret or Public
type Visibility uint8

const (
	Unset Visibility = iota
	Internal
	Secret
	Public
	Virtual
)

// Hint represents a solver hint
// it enables the solver to compute a Wire with a function provided at solving time
// using pre-defined inputs
type Hint struct {
	WireID int                // resulting wire ID to compute
	ID     hint.ID            // hint function id
	Inputs []LinearExpression // terms to inject in the hint function
}

func (r1c *R1C) String(coeffs []big.Int) string {
	var sbb strings.Builder
	sbb.WriteString("L[")
	r1c.L.string(&sbb, coeffs)
	sbb.WriteString("] * R[")
	r1c.R.string(&sbb, coeffs)
	sbb.WriteString("] = O[")
	r1c.O.string(&sbb, coeffs)
	sbb.WriteString("]")

	return sbb.String()
}

func (l LinearExpression) string(sbb *strings.Builder, coeffs []big.Int) {
	for i := 0; i < len(l); i++ {
		l[i].string(sbb, coeffs)
		if i+1 < len(l) {
			sbb.WriteString(" + ")
		}
	}
}

func (t Term) string(sbb *strings.Builder, coeffs []big.Int) {
	sbb.WriteString(coeffs[t.CoeffID()].String())
	sbb.WriteString("*")
	switch t.VariableVisibility() {
	case Internal:
		sbb.WriteString("i")
	case Public:
		sbb.WriteString("p")
	case Secret:
		sbb.WriteString("s")
	case Virtual:
		sbb.WriteString("v")
	case Unset:
		sbb.WriteString("u")
	default:
		panic("not implemented")
	}
	sbb.WriteString(strconv.Itoa(t.VariableID()))
}