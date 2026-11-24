package shoppingcarts_test

import (
	"net/http"
	"testing"

	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	shoppingcartsclient "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated/client"
	"github.com/KochKevin/effective-spoon-v2/internal/testutil/testenv"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/stretchr/testify/require"
)

const module = "shoppingcart"

func login(t *testing.T, env *testenv.Env) users.User {
	//Create User
	user, err := env.SeedUser(t)
	require.NoError(t, err, "no error seeding a user")

	//Login User
	err = env.App.Services.Input.EnterRfid(t.Context(), string(user.Code))
	require.NoError(t, err, "no error at login")

	return user
}

func Test_CreateEmptyCart(t *testing.T) {

	t.Parallel()

	env := testenv.New(t)

	c, err := shoppingcartsclient.NewClientWithResponses(env.URL + "/api/")
	require.NoError(t, err, "no error when creating client")

	user := login(t, env)

	//Create Cart
	res, err := c.PostShoppingCartsCurrentWithResponse(t.Context())
	require.NoError(t, err, "no error at posting shopping cart")

	//Check return value
	require.Equal(t, http.StatusCreated, res.StatusCode(), "has wrong status code")
	require.NotNil(t, res.Body, "response body is not nil")

	require.Equal(t, string(shoppingcarts.ShoppingCartActive), res.JSON201.Status, "cart should have correct status")
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

	user := login(t, env)

	//Create Cart
	_, err = c.PostShoppingCartsCurrentWithResponse(t.Context())
	require.NoError(t, err, "no error at posting shopping cart")

	//Create Product
	p, err := env.SeedProduct(t)
	require.NoError(t, err, "no error at seeding a product")

	//Add Product
	res, err := c.PostShoppingCartsCurrentIncreaseWithResponse(
		t.Context(),
		&shoppingcartsclient.PostShoppingCartsCurrentIncreaseParams{
			ProductID: p.Id.String(),
		},
	)
	require.NoError(t, err, "no error at posting shopping cart")

	//Check return value
	require.Equal(t, http.StatusOK, res.StatusCode(), "has wrong status code")
	require.NotNil(t, res.Body, "response body is not nil")
	require.Equal(t, string(shoppingcarts.ShoppingCartActive), res.JSON200.Status, "cart should have correct status")
	require.NotEmpty(t, res.JSON200.LineItems, "line items should be not empty")

	require.Len(t, res.JSON200.LineItems, 1, "only one line item")
	require.Equal(t, p.Id.String(), res.JSON200.LineItems[0].ProductId, "should be the correct product id")
	require.Equal(t, 1, res.JSON200.LineItems[0].Amount, "only one product")
	require.Equal(t, p.Name, res.JSON200.LineItems[0].ProductName, "should be same name")
	require.Equal(t, p.Price.GetAsEuro(), res.JSON200.LineItems[0].Price, "should be same price")

	require.Equal(t, p.Price.GetAsEuro(), res.JSON200.FullPrice, "full price should be only one product price")

	require.False(t, res.JSON200.UseEvent, "cart shouldn't use an event")
	require.Equal(t, user.Id.String(), res.JSON200.UserId, "cart must have the same user id as the user who created it")

}
