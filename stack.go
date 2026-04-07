package directives

// stack manages the context scope stack during parsing.
type stack struct {
	frames []ContextFrame
}

// push adds a frame to the stack.
func (s *stack) push(f ContextFrame) {
	s.frames = append(s.frames, f)
}

// pop removes and returns the top frame. Returns zero value and false if empty.
func (s *stack) pop() (ContextFrame, bool) {
	if len(s.frames) == 0 {
		return ContextFrame{}, false
	}
	f := s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-1]
	return f, true
}

// peek returns the top frame without removing it. Returns zero value and false if empty.
func (s *stack) peek() (ContextFrame, bool) {
	if len(s.frames) == 0 {
		return ContextFrame{}, false
	}
	return s.frames[len(s.frames)-1], true
}

// depth returns the current stack depth.
func (s *stack) depth() int {
	return len(s.frames)
}

// empty returns true if the stack has no frames.
func (s *stack) empty() bool {
	return len(s.frames) == 0
}

// contextStart returns the line number of the outermost context frame,
// or 0 if no context is active.
func (s *stack) contextStart() int {
	for _, f := range s.frames {
		if f.Kind == "context" {
			return f.Start
		}
	}
	return 0
}

// currentStart returns the start line of the innermost frame (zoom or context).
// Returns 0 if the stack is empty.
func (s *stack) currentStart() int {
	if len(s.frames) == 0 {
		return 0
	}
	return s.frames[len(s.frames)-1].Start
}

// popZooms pops all zoom frames until the stack is empty or a context frame is on top.
// Returns the number of frames popped.
func (s *stack) popZooms() int {
	count := 0
	for len(s.frames) > 0 && s.frames[len(s.frames)-1].Kind == "zoom" {
		s.frames = s.frames[:len(s.frames)-1]
		count++
	}
	return count
}

// clear removes all frames from the stack. Returns the number of frames removed.
func (s *stack) clear() int {
	n := len(s.frames)
	s.frames = nil
	return n
}
