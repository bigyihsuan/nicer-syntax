package evaluator

type Evaluable interface {
	Evaluate() NicerValue
}

type NicerValue struct {
	Value interface{}
}
