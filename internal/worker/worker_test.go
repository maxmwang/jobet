package worker

import (
	"context"
	"testing"

	"github.com/maxmwang/jobet/internal/mocks"
	"github.com/maxmwang/jobet/internal/proto"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWorker_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	mockPublisherA := mocks.NewMockPublisher(ctrl)
	mockPublisherB := mocks.NewMockPublisher(ctrl)
	mockPublisherC := mocks.NewMockPublisher(ctrl)
	mockWorker := Worker{
		publishers: []Publisher{
			mockPublisherA,
			mockPublisherB,
			mockPublisherC,
		},
	}

	mockBatch := &proto.ScrapeBatch{
		Priority: 1,
		Jobs: []*proto.ScrapeBatch_Job{
			{Company: "company", Title: "test"},
		},
	}

	mockPublisherA.EXPECT().Publish(ctx, mockBatch)
	mockPublisherB.EXPECT().Publish(ctx, mockBatch)
	mockPublisherC.EXPECT().Publish(ctx, mockBatch)

	err := mockWorker.publish(ctx, mockBatch)
	require.NoError(t, err)
}
