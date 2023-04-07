package main

type dwfa struct {
	states      int
	symbols     int
	startState  wfaState
	transitions map[wfaState]map[symbol]wfaTransition
}
type wfaState int
type wfaTransition struct {
	wfaState
	weight
}
type weight int
type specialSets struct {
	nonNegative map[wfaState]struct{}
	nonPositive map[wfaState]struct{}
}

type turingMachine struct {
	states      int
	symbols     int
	transitions map[tmState]map[symbol]tmTransition
}
type tmTransition struct {
	symbol
	direction
	tmState
}
type symbol int
type direction bool

const L direction = true
const LEFT direction = true
const R direction = false
const RIGHT direction = false

func (d direction) String() string {
	if d {
		return "L"
	}
	return "R"
}

type tmState int

const A tmState = 0
const B tmState = 1
const C tmState = 2
const D tmState = 3
const E tmState = 4
const F tmState = 5
const Z tmState = -1

const TMSTARTSTATE tmState = 0
const TMSTARTSYMBOL symbol = 0

const HALTSTATESTRING = "[HALT]"

func (tms tmState) String() string {
	if tms < 0 {
		return HALTSTATESTRING
	}
	return string(byte(tms) + 'A')
}

type acceptSet map[config]condition

type config struct {
	tmState    tmState
	tmSymbol   symbol
	leftState  wfaState
	rightState wfaState
}

type condition struct {
	mode
	weight
}

type mode int

const EQUAL mode = 0
const MOREOREQUAL mode = 1
const LESSOREQUAL mode = 2
const ACCEPTALL mode = 4

func MITMWFARverifier(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	return verifyValidTM(tm) &&
		verifyDeterministicWFA(leftWFA) &&
		verifyDeterministicWFA(rightWFA) &&
		verifySymbolCompatibility(tm, leftWFA, rightWFA) &&
		verifyLeadingBlankInvariant(leftWFA) &&
		verifyLeadingBlankInvariant(rightWFA) &&
		verifySpecialSets(leftWFA, leftSpecialSets) &&
		verifySpecialSets(rightWFA, rightSpecialSets) &&
		verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) &&
		verifyNoHaltingConfigAccepted(tm, acceptSet) &&
		verifyForwardClosed(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet)
}

func verifyValidTM(tm turingMachine) bool {
	if tm.states <= 0 || tm.symbols <= 0 {
		return false
	}
	return true
}

func verifyDeterministicWFA(wfa dwfa) bool {
	if wfa.states <= 0 || wfa.symbols <= 0 {
		return false
	}
	if wfa.startState < 0 || int(wfa.startState) >= wfa.states {
		return false
	}
	for i := 0; i < wfa.states; i++ {
		for j := 0; j < wfa.symbols; j++ {
			transition, ok := wfa.transitions[wfaState(i)][symbol(j)]
			if !ok {
				return false
			}
			if transition.wfaState < 0 || int(transition.wfaState) >= wfa.states {
				return false
			}
		}
	}
	return true
}

func verifySymbolCompatibility(tm turingMachine, leftWFA, rightWFA dwfa) bool {
	return tm.symbols == leftWFA.symbols && tm.symbols == rightWFA.symbols
}

func verifyLeadingBlankInvariant(wfa dwfa) bool {
	state := wfa.startState
	transition := wfa.transitions[state][0]
	return transition.wfaState == state && transition.weight == 0
}

func verifySpecialSets(wfa dwfa, specialSets specialSets) bool {
	for i := 0; i < wfa.states; i++ {
		for j := 0; j < wfa.symbols; j++ {
			transition := wfa.transitions[wfaState(i)][symbol(j)]
			if !transitionRetainsSpecialSets(wfaState(i), transition.wfaState, transition.weight, specialSets) {
				return false
			}
		}
	}
	return true
}

func transitionRetainsSpecialSets(startState, endState wfaState, weight weight, specialSets specialSets) bool {

	_, endNonPositive := specialSets.nonPositive[endState]
	_, startNonPositive := specialSets.nonPositive[startState]
	if endNonPositive {
		if !startNonPositive {
			return false
		}
		if weight > 0 {
			return false
		}
	}

	_, endNonNegative := specialSets.nonNegative[endState]
	_, startNonNegative := specialSets.nonNegative[startState]
	if endNonNegative {
		if !startNonNegative {
			return false
		}
		if weight < 0 {
			return false
		}
	}

	return true
}

func verifyStartConfigAccept(leftWFA, rightWFA dwfa, acceptSet acceptSet) bool {
	condition, ok := acceptSet[config{TMSTARTSTATE, TMSTARTSYMBOL, leftWFA.startState, rightWFA.startState}]
	if !ok {
		return false
	}
	switch condition.mode {
	case EQUAL:
		return condition.weight == 0
	case LESSOREQUAL:
		return condition.weight >= 0
	case MOREOREQUAL:
		return condition.weight <= 0
	case ACCEPTALL:
		return true
	}
	return false
}

func verifyNoHaltingConfigAccepted(tm turingMachine, acceptSet acceptSet) bool {
	for condition := range acceptSet {
		if condition.tmState < 0 || int(condition.tmState) >= tm.states {
			return false
		}
		if haltsNextStep(tm, condition.tmState, condition.tmSymbol) {
			return false
		}
	}
	return true
}

func haltsNextStep(tm turingMachine, tmState tmState, symbol symbol) bool {
	transition, ok := tm.transitions[tmState][symbol]
	if !ok {
		return true
	}
	if transition.tmState < 0 || int(transition.tmState) >= tm.states {
		return true
	}
	return false
}

func verifyForwardClosed(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	panic("unimplemented")
}