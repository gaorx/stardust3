package sdtaskcron

import (
	"github.com/gaorx/stardust3/sdjson"
)

type Action func(taskId, actionId string, args sdjson.Object) error
