package caskin_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestNewObjectDirectory(t *testing.T) {
	// object1 := &example.Object{ID: 1, ParentID: 0}
	object2 := &example.Object{ID: 2, ParentID: 1}
	// object3 := &example.Object{ID: 3, ParentID: 1}
	object4 := &example.Object{ID: 4, ParentID: 2}
	object5 := &example.Object{ID: 5, ParentID: 2}
	object6 := &example.Object{ID: 6, ParentID: 3}
	object7 := &example.Object{ID: 7, ParentID: 3}
	object8 := &example.Object{ID: 8, ParentID: 5}
	directory := []*caskin.Directory{
		{Object: object2, TopItemCount: 2},
		{Object: object4, TopItemCount: 4},
		{Object: object5, TopItemCount: 5},
		{Object: object6, TopItemCount: 6},
		{Object: object7, TopItemCount: 7},
		{Object: object8, TopItemCount: 8},
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(directory), func(i, j int) { directory[i], directory[j] = directory[j], directory[i] })
	od := caskin.NewObjectDirectory(directory)
	assert.Equal(t, uint64(3), od.Node[2].AllDirectoryCount)
	assert.Equal(t, uint64(19), od.Node[2].AllItemCount)
	assert.Equal(t, uint64(1), od.Node[5].AllDirectoryCount)
	assert.Equal(t, uint64(13), od.Node[5].AllItemCount)
}
