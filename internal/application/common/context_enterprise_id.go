package common

import "context"

const enterpriseIDKey = "Ctx_enterpriseID"

func SetEnterpriseID(ctx context.Context, enterpriseID string) context.Context {
	return context.WithValue(ctx, enterpriseIDKey, enterpriseID)
}

func GetEnterpriseID(ctx context.Context) (string, bool) {
	enterpriseID, ok := ctx.Value(enterpriseIDKey).(string)

	return enterpriseID, ok
}

func ExistsEnterpriseID(ctx context.Context) bool {
	_, ok := ctx.Value(enterpriseIDKey).(string)

	return ok
}
