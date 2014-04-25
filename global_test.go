package tardis

import (
	"gopkg.in/check.v1"
)

type GlobalSuite struct{}
var _ = check.Suite(&GlobalSuite{})

func (s *GlobalSuite) TestEpoch(c *check.C) {
    timestamp, err := Epoch("2014-04-03T20:39:54+00:00")
    c.Assert(err, check.IsNil)
    c.Assert(timestamp, check.Equals, int64(1396557594))

    timestamp, err = Epoch("2014-04-25T02:05:27Z")
    c.Assert(err, check.IsNil)
    c.Assert(timestamp, check.Equals, int64(1398391527))
}

func (s *GlobalSuite) TestParseEvent(c *check.C) {
	event, err := ParseEvent([]byte(`{"shop_id": 12345}`))
	c.Assert(err, check.IsNil)
	c.Assert(event["shop_id"], check.Equals, "12345")
}
