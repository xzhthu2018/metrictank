package expr

import (
	"fmt"
	"reflect"

	"github.com/raintank/metrictank/api/models"
	"gopkg.in/raintank/schema.v1"
)

type FuncSumSeries struct {
}

func NewSumSeries() Func {
	return FuncSumSeries{}
}

func (s FuncSumSeries) Signature() ([]argType, []argType) {
	return []argType{seriesList}, []argType{series}
}

func (s FuncSumSeries) Init(args []*expr) error {
	return nil
}

func (s FuncSumSeries) Depends(from, to uint32) (uint32, uint32) {
	return from, to
}

func (s FuncSumSeries) Exec(cache map[Req][]models.Series, in ...interface{}) ([]interface{}, error) {
	series, ok := in[0].([]models.Series)
	if !ok {
		return nil, ErrBadArgument{reflect.TypeOf([]models.Series{}), reflect.TypeOf(in[0])}
	}
	if len(series) == 1 {
		return []interface{}{series[0]}, nil
	}
	out := pointSlicePool.Get().([]schema.Point)
	for i := 0; i < len(series[0].Datapoints); i++ {
		point := schema.Point{
			Ts:  series[0].Datapoints[i].Ts,
			Val: 0,
		}
		for j := 0; j < len(series); j++ {
			point.Val += series[j].Datapoints[i].Val
		}
		out = append(out, point)
	}
	output := models.Series{
		Target:     fmt.Sprintf("sumSeries(%s)", "foo"),
		Datapoints: out,
		Interval:   series[0].Interval,
	}
	cache[Req{}] = append(cache[Req{}], output)
	return []interface{}{output}, nil
}
