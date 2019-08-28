package utils

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func StartXraySegment(ctx context.Context, operation string) (context.Context, *xray.Segment) {
	if xray.GetSegment(ctx) != nil {
		return xray.BeginSubsegment(ctx, operation)
	}

	return xray.BeginSegment(ctx, operation)
}
