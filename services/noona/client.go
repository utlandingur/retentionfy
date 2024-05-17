package noona

import (
	"context"

	"github.com/noona-hq/app-template/utils"
	noona "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

type Client struct {
	cfg    Config
	Client *noona.ClientWithResponses
}

func (a Client) GetUser() (*noona.User, error) {
	userResponse, err := a.Client.GetUserWithResponse(context.Background(), &noona.GetUserParams{
		Expand: &noona.Expand{"companies"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error getting user")
	}

	if userResponse.StatusCode() != 200 {
		return nil, errors.New("Error getting user")
	}

	if userResponse.JSON200 == nil {
		return nil, errors.New("Error getting user")
	}

	return userResponse.JSON200, nil
}

func (a Client) SetupWebhook(companyID string) error {
	webhook := noona.Webhook{
		Title:       utils.StringPtr("Example Webhook"),
		Description: utils.StringPtr("Watches event creation to do something valuable at some point. Maybe."),
		CallbackUrl: utils.StringPtr(a.cfg.AppBaseURL + "/webhook"),
		Company: func() *noona.ExpandableCompany {
			company := noona.ExpandableCompany{}
			company.FromID(noona.ID(companyID))
			return &company
		}(),
		Enabled: utils.BoolPtr(true),
		Headers: &noona.WebhookHeaders{
			{
				Key:    utils.StringPtr("Authorization"),
				Values: &[]string{"Bearer " + a.cfg.AppWebhookToken},
			},
		},
		Events: &noona.WebhookEvents{
			noona.WebhookEventEventCreated,
		},
	}

	webhookResponse, err := a.Client.CreateWebhookWithResponse(context.Background(), &noona.CreateWebhookParams{}, noona.CreateWebhookJSONRequestBody(webhook))
	if err != nil {
		return errors.Wrap(err, "Error creating webhook")
	}

	if webhookResponse.StatusCode() != 200 {
		return errors.New("Error creating webhook")
	}

	return nil
}

func (a Client) SetupSomeResource(companyID string) error {
	group := noona.CustomerGroup{
		Title:       utils.StringPtr("Some Title"),
		Description: utils.StringPtr("Some Description."),
		Company:     &companyID,
	}

	groupResponse, err := a.Client.CreateCustomerGroupWithResponse(context.Background(), &noona.CreateCustomerGroupParams{}, noona.CreateCustomerGroupJSONRequestBody(group))
	if err != nil {
		return errors.Wrap(err, "Error creating customer group")
	}

	if groupResponse.StatusCode() != 200 {
		return errors.New("Error creating customer group")
	}

	return nil
}
