package record

import "github.com/gohade/hade/framework"

type RecordProvider struct {
	framework.ServiceProvider
	c framework.Container
}
func (sp *RecordProvider) Name() string {
	return RecordKey
}
func (sp *RecordProvider) Register(c framework.Container) framework.NewInstance {
	return NewRecordService
}
func (sp *RecordProvider) IsDefer() bool {
	return false
}
func (sp *RecordProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
func (sp *RecordProvider) Boot(c framework.Container) error {
	return nil
}