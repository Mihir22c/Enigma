package main

import (
	"fmt"
	"strings"
)

// Rotor represents a single rotor in the Enigma machine
type Rotor struct {
	wiring   string // The rotor's internal wiring
	position int    // The current position of the rotor
	notch    int    // The notch position for rotor stepping
}

// Reflector represents the reflector in the Enigma machine
type Reflector struct {
	wiring string // The reflector's wiring
}

// Plugboard represents the plugboard for letter swapping
type Plugboard struct {
	swaps map[rune]rune
}

// NewRotor creates a new rotor with a given wiring and notch position
func NewRotor(wiring string, notch int) *Rotor {
	return &Rotor{
		wiring:   wiring,
		position: notch,
		notch:    notch,
	}
}

// NewReflector creates a new reflector with given wiring
func NewReflector(wiring string) *Reflector {
	return &Reflector{
		wiring: wiring,
	}
}

// NewPlugboard creates a new plugboard with given swaps
func NewPlugboard(swaps map[rune]rune) *Plugboard {
	return &Plugboard{
		swaps: swaps,
	}
}

// Step advances the rotor position and handles stepping logic
func (r *Rotor) Step() {
	r.position = (r.position + 1) % 26
}

// Encrypt processes a single character through the Enigma machine
func (e *Enigma) Encrypt(c rune) rune {
	// Plugboard swap
	if swapped, ok := e.plugboard.swaps[c]; ok {
		c = swapped
	}

	// Pass through rotors
	for _, rotor := range e.rotors {
		c = rotor.process(c)
	}

	// Pass through reflector
	c = e.reflector.reflect(c)

	// Pass back through rotors in reverse order
	for i := len(e.rotors) - 1; i >= 0; i-- {
		c = e.rotors[i].reverseProcess(c)
	}

	// Plugboard swap
	if swapped, ok := e.plugboard.swaps[c]; ok {
		c = swapped
	}

	// Rotate the first rotor, and handle stepping logic for the other rotors
	e.rotors[0].Step()
	if e.rotors[0].position == e.rotors[0].notch {
		for i := 1; i < len(e.rotors); i++ {
			e.rotors[i].Step()
			if e.rotors[i].position != e.rotors[i].notch {
				break
			}
		}
	}

	return c
}

// Process a character through a rotor
func (r *Rotor) process(c rune) rune {
	offset := r.position
	c = (c-'A'+rune(offset))%26 + 'A'
	return rune(r.wiring[c-'A'])
}

// Reverse process a character through a rotor
func (r *Rotor) reverseProcess(c rune) rune {
	offset := r.position
	c = (c-'A'-rune(offset)+26)%26 + 'A'
	index := strings.IndexRune(r.wiring, c)
	return rune('A' + index)
}

// Reflect a character through the reflector
func (r *Reflector) reflect(c rune) rune {
	return rune(r.wiring[c-'A'])
}

// Enigma represents the entire Enigma machine
type Enigma struct {
	rotors    []*Rotor
	reflector *Reflector
	plugboard *Plugboard
}

// NewEnigma creates a new Enigma machine with the given components
func NewEnigma(rotors []*Rotor, reflector *Reflector, plugboard *Plugboard) *Enigma {
	return &Enigma{
		rotors:    rotors,
		reflector: reflector,
		plugboard: plugboard,
	}
}

func main() {
	// Create rotors with example wirings and notch positions
	rotor1 := NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 17)
	rotor2 := NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", 5)
	rotor3 := NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", 22)

	// Create reflector with example wiring
	reflector := NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT")

	// Create plugboard with example swaps
	plugboard := NewPlugboard(map[rune]rune{
		'A': 'Z', 'Z': 'A',
		'B': 'Y', 'Y': 'B',
	})

	// Create the Enigma machine
	enigma := NewEnigma([]*Rotor{rotor1, rotor2, rotor3}, reflector, plugboard)

	// Encrypt a message
	message := "HELLOENIGMA"
	encryptedMessage := ""
	for _, char := range message {
		if char >= 'A' && char <= 'Z' {
			encryptedMessage += string(enigma.Encrypt(char))
		} else {
			encryptedMessage += string(char)
		}
	}

	fmt.Println("Encrypted message:", encryptedMessage)

	// Decrypt the message by reusing the Enigma machine (since it's symmetric)
	decryptedMessage := ""
	for _, char := range encryptedMessage {
		if char >= 'A' && char <= 'Z' {
			decryptedMessage += string(enigma.Encrypt(char))
		} else {
			decryptedMessage += string(char)
		}
	}

	fmt.Println("Decrypted message:", decryptedMessage)
}
