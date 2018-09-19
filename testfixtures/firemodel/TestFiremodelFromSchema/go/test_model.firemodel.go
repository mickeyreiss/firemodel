// DO NOT EDIT - Code generated by firemodel (dev).

package firemodel

import (
	firestore "cloud.google.com/go/firestore"
	"fmt"
	runtime "github.com/mickeyreiss/firemodel/runtime"
	latlng "google.golang.org/genproto/googleapis/type/latlng"
	"time"
)

// A Test is a test model.
//
// Firestore document location: /users/{user_id}/test_models/{test_model_id}
type TestModel struct {
	// The name.
	Name string `firestore:"name"`
	// The age.
	Age int64 `firestore:"age"`
	// The number pi.
	Pi float64 `firestore:"pi"`
	// The birth date.
	Birthdate time.Time `firestore:"birthdate"`
	// True if it is good.
	IsGood bool `firestore:"isGood"`
	// TODO: Add comment to TestModel.data.
	Data []byte `firestore:"data"`
	// TODO: Add comment to TestModel.friend.
	Friend *firestore.DocumentRef `firestore:"friend"`
	// TODO: Add comment to TestModel.location.
	Location *latlng.LatLng `firestore:"location"`
	// TODO: Add comment to TestModel.colors.
	Colors []string `firestore:"colors"`
	// TODO: Add comment to TestModel.directions.
	Directions []TestDirection `firestore:"directions"`
	// TODO: Add comment to TestModel.models.
	Models []*TestStruct `firestore:"models"`
	// TODO: Add comment to TestModel.refs.
	Refs []*firestore.DocumentRef `firestore:"refs"`
	// TODO: Add comment to TestModel.model_refs.
	ModelRefs []*firestore.DocumentRef `firestore:"modelRefs"`
	// TODO: Add comment to TestModel.meta.
	Meta map[string]interface{} `firestore:"meta"`
	// TODO: Add comment to TestModel.meta_strs.
	MetaStrs map[string]string `firestore:"metaStrs"`
	// TODO: Add comment to TestModel.direction.
	Direction TestDirection `firestore:"direction"`
	// TODO: Add comment to TestModel.test_file.
	TestFile *runtime.File `firestore:"testFile"`
	// TODO: Add comment to TestModel.url.
	Url runtime.URL `firestore:"url"`
	// TODO: Add comment to TestModel.nested.
	Nested *TestStruct `firestore:"nested"`

	// Creation timestamp.
	CreatedAt time.Time `firestore:"createdAt,serverTimestamp"`
	// Update timestamp.
	UpdatedAt time.Time `firestore:"updatedAt,serverTimestamp"`
}

// TestModelPath returns the path to a particular TestModel in Firestore.
func TestModelPath(userId string, testModelId string) string {
	return fmt.Sprintf("users/%s/test_models/%s", userId, testModelId)
}
