package sdtaskcron

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/robfig/cron/v3"
)

var (
	cronParser = cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
)

func parseSched(s string) (cron.Schedule, error) {
	sch, err := cronParser.Parse(s)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return sch, nil
}
