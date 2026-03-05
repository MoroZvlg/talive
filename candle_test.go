package talive_test

type testCandle struct {
	open, high, low, close, volume float64
}

func (c *testCandle) Open() float64   { return c.open }
func (c *testCandle) High() float64   { return c.high }
func (c *testCandle) Low() float64    { return c.low }
func (c *testCandle) Close() float64  { return c.close }
func (c *testCandle) Volume() float64 { return c.volume }
