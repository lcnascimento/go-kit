package temporalclient

import "context"

func onWarnSearchAttributes(ctx context.Context, namespace string, err error) {
	if err == nil {
		return
	}

	logger.WarnContext(ctx, "unable to ensure search attributes", "namespace", namespace, "error", err.Error())
}
