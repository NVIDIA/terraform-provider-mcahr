// SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package systemdefer

import (
	"testing"
)

func TestNewDeferFunctionList(t *testing.T) {
	t.Parallel()

	// when
	deferList := NewDeferFunctionList()

	// then
	if deferList == nil {
		t.Error("NewDeferFunctionList should not return nil")
	}

	if deferList.deferList == nil {
		t.Error("deferList.deferList should not be nil")
	}

	if deferList.Size() != 0 {
		t.Errorf("expected empty defer list, got %d items", deferList.Size())
	}
}

func TestAddDefer_SingleFunction(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	called := false
	testFunc := func() {
		called = true
	}

	// when
	deferList.AddDefer(testFunc)

	// then
	if deferList.Size() != 1 {
		t.Errorf("expected 1 defer function, got %d", deferList.Size())
	}

	// Function should not be called yet
	if called {
		t.Error("defer function should not be called until ExecuteAll is called")
	}
}

func TestAddDefer_MultipleFunctions(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	var calls []int

	func1 := func() { calls = append(calls, 1) }
	func2 := func() { calls = append(calls, 2) }
	func3 := func() { calls = append(calls, 3) }

	// when
	deferList.AddDefer(func1)
	deferList.AddDefer(func2)
	deferList.AddDefer(func3)

	// then
	if deferList.Size() != 3 {
		t.Errorf("expected 3 defer functions, got %d", deferList.Size())
	}

	// Functions should not be called yet
	if len(calls) != 0 {
		t.Error("defer functions should not be called until ExecuteAll is called")
	}
}

func TestExecuteAll_SingleFunction(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	called := false
	testFunc := func() {
		called = true
	}
	deferList.AddDefer(testFunc)

	// when
	deferList.ExecuteAll()

	// then
	if !called {
		t.Error("defer function should have been called")
	}
}

func TestExecuteAll_MultipleFunctions(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	var calls []int

	func1 := func() { calls = append(calls, 1) }
	func2 := func() { calls = append(calls, 2) }
	func3 := func() { calls = append(calls, 3) }

	deferList.AddDefer(func1)
	deferList.AddDefer(func2)
	deferList.AddDefer(func3)

	// when
	deferList.ExecuteAll()

	// then
	expectedCalls := []int{1, 2, 3}
	if len(calls) != len(expectedCalls) {
		t.Errorf("expected %d calls, got %d", len(expectedCalls), len(calls))
	}

	for i, expected := range expectedCalls {
		if i >= len(calls) || calls[i] != expected {
			t.Errorf("expected call sequence %v, got %v", expectedCalls, calls)
			break
		}
	}
}

func TestExecuteAll_OrderOfExecution(t *testing.T) {
	t.Parallel()

	// given - test that functions are called in the order they were added (FIFO, not LIFO like Go's defer)
	deferList := NewDeferFunctionList()
	var executionOrder []string

	first := func() { executionOrder = append(executionOrder, "first") }
	second := func() { executionOrder = append(executionOrder, "second") }
	third := func() { executionOrder = append(executionOrder, "third") }

	deferList.AddDefer(first)
	deferList.AddDefer(second)
	deferList.AddDefer(third)

	// when
	deferList.ExecuteAll()

	// then
	expectedOrder := []string{"first", "second", "third"}
	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(executionOrder) || executionOrder[i] != expected {
			t.Errorf("expected execution order %v, got %v", expectedOrder, executionOrder)
			break
		}
	}
}

func TestExecuteAll_EmptyList(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()

	// when - should not panic
	deferList.ExecuteAll()

	// then - test passes if no panic occurred
}

func TestExecuteAll_WithPanic(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	var calls []int

	func1 := func() { calls = append(calls, 1) }
	panicFunc := func() { panic("test panic") }
	func3 := func() { calls = append(calls, 3) }

	deferList.AddDefer(func1)
	deferList.AddDefer(panicFunc)
	deferList.AddDefer(func3)

	// when/then - should panic, but we can verify the first function was called
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from defer function")
		}
		// Verify first function was called before panic
		if len(calls) != 1 || calls[0] != 1 {
			t.Errorf("expected first function to be called before panic, got calls: %v", calls)
		}
	}()

	deferList.ExecuteAll()
}

func TestExecuteAll_CalledMultipleTimes(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	callCount := 0
	testFunc := func() {
		callCount++
	}
	deferList.AddDefer(testFunc)

	// when
	deferList.ExecuteAll()
	deferList.ExecuteAll() // Call again

	// then - function should be called each time ExecuteAll is invoked
	if callCount != 2 {
		t.Errorf("expected function to be called 2 times, got %d", callCount)
	}
}

func TestAddDefer_AfterExecuteAll(t *testing.T) {
	t.Parallel()

	// given
	deferList := NewDeferFunctionList()
	firstCallCount := 0
	secondCallCount := 0

	firstFunc := func() { firstCallCount++ }
	secondFunc := func() { secondCallCount++ }

	deferList.AddDefer(firstFunc)
	deferList.ExecuteAll()

	// when - add another function after ExecuteAll
	deferList.AddDefer(secondFunc)
	deferList.ExecuteAll()

	// then
	if firstCallCount != 2 {
		t.Errorf("expected first function to be called 2 times, got %d", firstCallCount)
	}
	if secondCallCount != 1 {
		t.Errorf("expected second function to be called 1 time, got %d", secondCallCount)
	}
}
