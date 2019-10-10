package saga

var _ Saga = (*Workflow)(nil) // ensure Workflow implements Saga
