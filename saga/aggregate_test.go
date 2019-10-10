package saga

var _ Saga = (*Aggregate)(nil) // ensure Aggregate implements Saga
