package domain

type Clock interface{ Now() int64 }

type FixedClock struct{ Value int64 }

func (c FixedClock) Now() int64 { return c.Value }

func Advance(clock FixedClock, delta int64) FixedClock { return FixedClock{Value: clock.Value + delta} }
