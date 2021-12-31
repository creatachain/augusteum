package pubsub_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/creatachain/augusteum/libs/log"

	"github.com/creatachain/augusteum/libs/pubsub"
	"github.com/creatachain/augusteum/libs/pubsub/query"
)

func TestExample(t *testing.T) {
	s := pubsub.NewServer()
	s.SetLogger(log.TestingLogger())
	err := s.Start()
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := s.Stop(); err != nil {
			t.Error(err)
		}
	})

	ctx := context.Background()
	subscription, err := s.Subscribe(ctx, "example-client", query.MustParse("msm.account.name='John'"))
	require.NoError(t, err)
	err = s.PublishWithEvents(ctx, "Tombstone", map[string][]string{"msm.account.name": {"John"}})
	require.NoError(t, err)
	assertReceive(t, "Tombstone", subscription.Out())
}
