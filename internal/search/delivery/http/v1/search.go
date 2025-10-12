package http

import (
	"context"

	gen "github.com/keenywheels/go-spy/internal/api/v1"
	"github.com/keenywheels/go-spy/internal/search/delivery/http/security"
	"github.com/keenywheels/go-spy/pkg/ctxutils"
	"github.com/keenywheels/go-spy/pkg/httputils"
)

// StartSearch TODO: заменить мок
func (c *Controller) StartSearch(
	ctx context.Context,
	req *gen.StartSearchRequest,
	params gen.StartSearchParams,
) (gen.StartSearchRes, error) {
	op := "Controller.GetAllInterest"
	log := ctxutils.GetLogger(ctx)

	// validate that client using his token
	client := security.GetClientFromContext(ctx)
	if client != params.XClient {
		log.Errorf("[%s] client %s is using token for client %s", op, params.XClient, client)

		return &gen.StartSearchForbidden{
			Error: httputils.ErrorForbidden,
		}, nil
	}

	// validate request
	if err := req.Validate(); err != nil {
		log.Errorf("[%s] failed to validate request: %v", op, err)

		return &gen.StartSearchBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	log.Infof("%s: got start search request for search: %+v", op, req)

	resp := gen.StartSearchOKApplicationJSON(getMockData())

	return &resp, nil
}

// getMockData simply returns some mock data
func getMockData() []gen.SearchMessage {
	return []gen.SearchMessage{
		{Message: "Boy favourable day can introduced sentiments entreaties. Noisier carried of in warrant because. So mr plate seems cause chief widen first. Two differed husbands met screened his. Bed was form wife out ask draw. Wholly coming at we no enable. Offending sir delivered questions now new met. Acceptance she interested new boisterous day discretion celebrated."},
		{Message: "Spot of come to ever hand as lady meet on. Delicate contempt received two yet advanced. Gentleman as belonging he commanded believing dejection in by. On no am winding chicken so behaved. Its preserved sex enjoyment new way behaviour. Him yet devonshire celebrated especially. Unfeeling one provision are smallness resembled repulsive."},
		{Message: "Manor we shall merit by chief wound no or would. Oh towards between subject passage sending mention or it. Sight happy do burst fruit to woody begin at. Assurance perpetual he in oh determine as. The year paid met him does eyes same. Own marianne improved sociable not out. Thing do sight blush mr an. Celebrated am announcing delightful remarkably we in literature it solicitude. Design use say piqued any gay supply. Front sex match vexed her those great."},
		{Message: "Started earnest brother believe an exposed so. Me he believing daughters if forfeited at furniture. Age again and stuff downs spoke. Late hour new nay able fat each sell. Nor themselves age introduced frequently use unsatiable devonshire get. They why quit gay cold rose deal park. One same they four did ask busy. Reserved opinions fat him nay position. Breakfast as zealously incommode do agreeable furniture. One too nay led fanny allow plate."},
		{Message: "Examine she brother prudent add day ham. Far stairs now coming bed oppose hunted become his. You zealously departure had procuring suspicion. Books whose front would purse if be do decay. Quitting you way formerly disposed perceive ladyship are. Common turned boy direct and yet."},
	}
}
