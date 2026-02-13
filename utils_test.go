package main

import (
    "testing"
)

func TestIdentifyAgentCommitterEmail(t *testing.T) {
    name, exists := IdentifyAgentCommitterEmail("209825114+claude[bot]@users.noreply.github.com")
    if name != "Claude" || exists != true {
        t.Errorf(`IdentifyAgentCommitterEmail() for claude address = %q, %v, expecting %#q`, name, exists, "Claude")
    }
}

func TestIdentifyAgentCommitterEmailNone(t *testing.T) {
    name, exists := IdentifyAgentCommitterEmail("user@example.com")
    if name != "" || exists != false {
        t.Errorf(`IdentifyAgentCommitterEmail() for non agent address = %q, %v, expecting nothing to exist`, name, exists)
    }
}

