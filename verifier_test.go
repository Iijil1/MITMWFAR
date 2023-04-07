package main

import (
	"fmt"
	"testing"
)

func TestTypes(t *testing.T) {
	t.Run("TmStateString", func(t *testing.T) {
		var state tmState = 1
		testString := fmt.Sprint(state)
		if testString != "B" {
			t.Fail()
		}
		var halt tmState = -1
		haltString := fmt.Sprint(halt)
		if haltString != HALTSTATESTRING {
			t.Fail()
		}
	})
	t.Run("DirectionString", func(t *testing.T) {
		var left direction = L
		leftString := fmt.Sprint(left)
		if leftString != "L" {
			t.Fail()
		}
		var right direction = RIGHT
		rightString := fmt.Sprint(right)
		if rightString != "R" {
			t.Fail()
		}
	})
}

func TestVerifyValidTM(t *testing.T) {
	t.Run("NoStates", func(t *testing.T) {
		tm := turingMachine{
			states:      0,
			symbols:     2,
			transitions: map[tmState]map[symbol]tmTransition{},
		}
		if verifyValidTM(tm) {
			t.Fail()
		}
	})
	t.Run("NoSymbols", func(t *testing.T) {
		tm := turingMachine{
			states:      2,
			symbols:     0,
			transitions: map[tmState]map[symbol]tmTransition{},
		}
		if verifyValidTM(tm) {
			t.Fail()
		}
	})
	t.Run("CorrectTM", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A},
					1: {1, R, Z}},
			},
		}
		if !verifyValidTM(tm) {
			t.Fail()
		}
	})
}

func TestVerifyDeterministic(t *testing.T) {
	t.Run("NoStates", func(t *testing.T) {
		wfa := dwfa{
			states:      0,
			symbols:     1,
			startState:  0,
			transitions: map[wfaState]map[symbol]wfaTransition{},
		}
		if verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
	t.Run("NoSymbols", func(t *testing.T) {
		wfa := dwfa{
			states:      1,
			symbols:     0,
			startState:  0,
			transitions: map[wfaState]map[symbol]wfaTransition{},
		}
		if verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
	t.Run("OutOfBoundStart", func(t *testing.T) {
		wfa := dwfa{
			states:      1,
			symbols:     1,
			startState:  1,
			transitions: map[wfaState]map[symbol]wfaTransition{},
		}
		if verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
	t.Run("IncompleteTransitionStateMap", func(t *testing.T) {
		wfa := dwfa{
			states:      2,
			symbols:     2,
			startState:  0,
			transitions: map[wfaState]map[symbol]wfaTransition{0: {0: {0, 0}, 1: {0, 0}}},
		}
		if verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
	t.Run("IncompleteTransitionSymbolMap", func(t *testing.T) {
		wfa := dwfa{
			states:     2,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {0, 0}},
				1: {0: {0, 0}}},
		}
		if verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
	t.Run("OutOfBoundTransition", func(t *testing.T) {
		wfa := dwfa{
			states:     2,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 0}},
				1: {0: {0, 0},
					1: {2, 0}}},
		}
		if verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
	t.Run("CorrectWFA", func(t *testing.T) {
		wfa := dwfa{
			states:     2,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 1},
					1: {1, 0}},
				1: {0: {1, 2},
					1: {0, -2}}},
		}
		if !verifyDeterministicWFA(wfa) {
			t.Fail()
		}
	})
}

func TestVerifySymbolCompatibility(t *testing.T) {
	t.Run("DifferentLeft", func(t *testing.T) {
		tm := turingMachine{symbols: 2}
		leftWFA := dwfa{symbols: 3}
		rightWFA := dwfa{symbols: 2}
		if verifySymbolCompatibility(tm, leftWFA, rightWFA) {
			t.Fail()
		}
	})
	t.Run("DifferentRight", func(t *testing.T) {
		tm := turingMachine{symbols: 2}
		leftWFA := dwfa{symbols: 2}
		rightWFA := dwfa{symbols: 3}
		if verifySymbolCompatibility(tm, leftWFA, rightWFA) {
			t.Fail()
		}
	})
	t.Run("Correct", func(t *testing.T) {
		tm := turingMachine{symbols: 2}
		leftWFA := dwfa{symbols: 2}
		rightWFA := dwfa{symbols: 2}
		if !verifySymbolCompatibility(tm, leftWFA, rightWFA) {
			t.Fail()
		}
	})
}

func TestVerifyLeadingBlankInvariant(t *testing.T) {
	t.Run("WrongTransitionState", func(t *testing.T) {
		wfa := dwfa{
			states:      2,
			symbols:     2,
			startState:  0,
			transitions: map[wfaState]map[symbol]wfaTransition{0: {0: {1, 0}}},
		}
		if verifyLeadingBlankInvariant(wfa) {
			t.Fail()
		}
	})
	t.Run("WrongTransitionWeight", func(t *testing.T) {
		wfa := dwfa{
			states:      2,
			symbols:     2,
			startState:  0,
			transitions: map[wfaState]map[symbol]wfaTransition{0: {0: {0, 1}}},
		}
		if verifyLeadingBlankInvariant(wfa) {
			t.Fail()
		}
	})
	t.Run("WrongStartState", func(t *testing.T) {
		wfa := dwfa{
			states:      2,
			symbols:     2,
			startState:  1,
			transitions: map[wfaState]map[symbol]wfaTransition{0: {0: {0, 0}}},
		}
		if verifyLeadingBlankInvariant(wfa) {
			t.Fail()
		}
	})
	t.Run("CorrectTransition", func(t *testing.T) {
		wfa := dwfa{
			states:      2,
			symbols:     2,
			startState:  0,
			transitions: map[wfaState]map[symbol]wfaTransition{0: {0: {0, 0}}},
		}
		if !verifyLeadingBlankInvariant(wfa) {
			t.Fail()
		}
	})
	t.Run("CorrectTransitionAlternateStart", func(t *testing.T) {
		wfa := dwfa{
			states:      2,
			symbols:     2,
			startState:  1,
			transitions: map[wfaState]map[symbol]wfaTransition{1: {0: {1, 0}}},
		}
		if !verifyLeadingBlankInvariant(wfa) {
			t.Fail()
		}
	})
}

func TestVerifySpecialSets(t *testing.T) {
	t.Run("EmptySet", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{},
			nonPositive: map[wfaState]struct{}{},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 1}},
				1: {0: {2, -1},
					1: {1, 0}},
				2: {0: {2, 1},
					1: {3, -1}},
				3: {0: {3, 1},
					1: {3, 0}},
			},
		}
		if !verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
	t.Run("NoWeights", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{0: {}, 1: {}, 2: {}, 3: {}},
			nonPositive: map[wfaState]struct{}{0: {}, 1: {}, 2: {}, 3: {}},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 0}},
				1: {0: {2, 0},
					1: {1, 0}},
				2: {0: {2, 0},
					1: {3, 0}},
				3: {0: {3, 0},
					1: {3, 0}},
			},
		}
		if !verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
	t.Run("CorrectSets", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{0: {}, 1: {}},
			nonPositive: map[wfaState]struct{}{0: {}, 1: {}},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 0}},
				1: {0: {2, -1},
					1: {1, 0}},
				2: {0: {2, 1},
					1: {3, 0}},
				3: {0: {3, -1},
					1: {3, 0}},
			},
		}
		if !verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
	t.Run("InternalPositive", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{0: {}, 1: {}},
			nonPositive: map[wfaState]struct{}{0: {}, 1: {}},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 1}},
				1: {0: {2, -1},
					1: {1, 0}},
				2: {0: {2, 1},
					1: {3, 0}},
				3: {0: {3, -1},
					1: {3, 0}},
			},
		}
		if verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
	t.Run("InternalNegative", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{0: {}, 1: {}},
			nonPositive: map[wfaState]struct{}{0: {}, 1: {}},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, -1}},
				1: {0: {2, -1},
					1: {1, 0}},
				2: {0: {2, 1},
					1: {3, 0}},
				3: {0: {3, -1},
					1: {3, 0}},
			},
		}
		if verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
	t.Run("NonClosedPositive", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{0: {}, 1: {}},
			nonPositive: map[wfaState]struct{}{2: {}, 3: {}},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 0}},
				1: {0: {2, 0},
					1: {1, 0}},
				2: {0: {2, 0},
					1: {3, 0}},
				3: {0: {3, 0},
					1: {3, 0}},
			},
		}
		if verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
	t.Run("NonClosedNegative", func(t *testing.T) {
		specialSets := specialSets{
			nonNegative: map[wfaState]struct{}{2: {}, 3: {}},
			nonPositive: map[wfaState]struct{}{0: {}, 1: {}},
		}
		wfa := dwfa{
			states:     4,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 0}},
				1: {0: {2, 0},
					1: {1, 0}},
				2: {0: {2, 0},
					1: {3, 0}},
				3: {0: {3, 0},
					1: {3, 0}},
			},
		}
		if verifySpecialSets(wfa, specialSets) {
			t.Fail()
		}
	})
}

func TestVerifyStartConfigAccept(t *testing.T) {
	t.Run("MissingConfig", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{}
		if verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("FailedEqualCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {EQUAL, 1}}
		if verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("CorrectEqualCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {EQUAL, 0}}
		if !verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("FailedMoreCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {MOREOREQUAL, 1}}
		if verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("CorrectMoreCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {MOREOREQUAL, -1}}
		if !verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("CorrectMoreConditionAtEqual", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {MOREOREQUAL, 0}}
		if !verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("FailedLessCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {LESSOREQUAL, -1}}
		if verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("CorrectLessCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {LESSOREQUAL, 1}}
		if !verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("CorrectLessConditionAtEqual", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {LESSOREQUAL, 0}}
		if !verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
	t.Run("CorrectAllCondition", func(t *testing.T) {
		leftWFA := dwfa{startState: 0}
		rightWFA := dwfa{startState: 0}
		acceptSet := map[config]condition{{TMSTARTSTATE, TMSTARTSYMBOL, 0, 0}: {ACCEPTALL, 42}}
		if !verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) {
			t.Fail()
		}
	})
}

func TestVerifyNoHaltingConfigAccepted(t *testing.T) {
	t.Run("OutOfBoundAcceptConfig", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A},
					1: {1, R, Z}},
			},
		}
		acceptSet := map[config]condition{
			{C, 0, 0, 0}: {EQUAL, 0},
		}
		if verifyNoHaltingConfigAccepted(tm, acceptSet) {
			t.Fail()
		}
	})
	t.Run("AcceptHalting", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A},
					1: {1, R, Z}},
			},
		}
		acceptSet := map[config]condition{
			{B, 1, 0, 0}: {EQUAL, 0},
		}
		if verifyNoHaltingConfigAccepted(tm, acceptSet) {
			t.Fail()
		}
	})
	t.Run("AcceptUndef", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A}},
			},
		}
		acceptSet := map[config]condition{
			{B, 1, 0, 0}: {EQUAL, 0},
		}
		if verifyNoHaltingConfigAccepted(tm, acceptSet) {
			t.Fail()
		}
	})
	t.Run("AcceptCorrect", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A},
					1: {1, R, Z}},
			},
		}
		acceptSet := map[config]condition{
			{A, 0, 0, 0}: {EQUAL, 0},
			{A, 1, 0, 0}: {EQUAL, 0},
			{B, 0, 0, 0}: {EQUAL, 0},
		}
		if !verifyNoHaltingConfigAccepted(tm, acceptSet) {
			t.Fail()
		}
	})
}