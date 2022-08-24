package gocqlastra

import (
	"reflect"
	"unsafe"

	"github.com/gocql/gocql"
)

// Replace the release version of Cassandra with a version that is not 4.0.0 when using Astra. This avoids the
// `system.peers_v2` problem.
func replaceHostInfoVersionHorribleHackPleaseRemoveMe(host *gocql.HostInfo) {
	rs := reflect.ValueOf(host).Elem()
	rf := rs.FieldByName("version")
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()

	major := rf.FieldByName("Major")
	minor := rf.FieldByName("Minor")
	patch := rf.FieldByName("Patch")

	major.SetInt(3)
	minor.SetInt(11)
	patch.SetInt(0)
}

type HostFilter struct{}

func (h HostFilter) Accept(host *gocql.HostInfo) bool {
	replaceHostInfoVersionHorribleHackPleaseRemoveMe(host)
	return true
}
