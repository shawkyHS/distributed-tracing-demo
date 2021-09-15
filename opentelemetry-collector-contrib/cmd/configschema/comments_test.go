// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configschema

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFieldComments(t *testing.T) {
	v := reflect.ValueOf(testStruct{})
	comments, err := commentsForStruct(v, testDR("../.."))
	assert.NoError(t, err)
	assert.Equal(t, "embedded, package qualified comment\n", comments["Duration"])
	assert.Equal(t, "testStruct comment\n", comments["_struct"])
}

func TestExternalType(t *testing.T) {
	u, err := uuid.NewUUID()
	assert.NoError(t, err)
	v := reflect.ValueOf(u)
	comments, err := commentsForStruct(v, testDR("."))
	assert.NoError(t, err)
	assert.Equal(t, "A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC\n4122.\n", comments["_struct"])
}
