package shoppingcarts_test

import (
	"net/http"
	"testing"

	shoppingcartsclient "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated/client"
	"github.com/KochKevin/effective-spoon-v2/internal/testutil/testenv"
	"github.com/stretchr/testify/require"
)

const module = "shoppingcart"

func Test_CreateEmptyCart(t *testing.T) {

	t.Parallel()

	env := testenv.New(t)

	c, err := shoppingcartsclient.NewClientWithResponses(env.URL + "/api/")
	require.NoError(t, err, "no error when creating client")

	//Create User
	user, err := env.SeedUser(t)
	require.NoError(t, err, "no error seeding a user")

	//Login User
	err = env.App.Services.Input.EnterRfid(t.Context(), string(user.Code))
	require.NoError(t, err, "no error at login")

	//Create Cart
	res, err := c.PostShoppingCartsCurrentWithResponse(t.Context())
	require.NoError(t, err, "no error at posting shopping cart")

	//Check return value
	require.Equal(t, http.StatusCreated, res.StatusCode(), "has wrong status code")
	require.NotNil(t, res.Body, "response body is not nil")

	require.Empty(t, res.JSON201.LineItems, "there shouldn't be any line items")
	require.Zero(t, res.JSON201.FullPrice, "the price should be zero")
	require.False(t, res.JSON201.UseEvent, "cart shouldn't use an event")
	require.Equal(t, user.Id.String(), res.JSON201.UserId, "cart must have the same user id as the user who created it")

}

func Test_CreateCartAndAddProduct(t *testing.T) {

	t.Parallel()

	env := testenv.New(t)

	c, err := shoppingcartsclient.NewClientWithResponses(env.URL + "/api/")
	require.NoError(t, err, "no error when creating client")

	//Create User
	user, err := env.SeedUser(t)
	require.NoError(t, err, "no error seeding a user")

	//Login User
	err = env.App.Services.Input.EnterRfid(t.Context(), string(user.Code))
	require.NoError(t, err, "no error at login")

	//Create Cart
	res, err := c.PostShoppingCartsCurrentWithResponse(t.Context())
	require.NoError(t, err, "no error at posting shopping cart")

	//Add Product
	res, err := c.PostShoppingCartsCurrentIncrease(
		t.Context(),
		*shoppingcartsclient.PostShoppingCartsCurrentIncreaseParams{
			ProductID: ,
		},
	)
	require.NoError(t, err, "no error at posting shopping cart")

	//Check return value
	require.Equal(t, http.StatusCreated, res.StatusCode(), "has wrong status code")
	require.NotNil(t, res.Body, "response body is not nil")

	require.Empty(t, res.JSON201.LineItems, "there shouldn't be any line items")
	require.Zero(t, res.JSON201.FullPrice, "the price should be zero")
	require.False(t, res.JSON201.UseEvent, "cart shouldn't use an event")
	require.Equal(t, user.Id.String(), res.JSON201.UserId, "cart must have the same user id as the user who created it")

}
