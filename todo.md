- Auth System
    - Provide test endpoint to login with user code
    - Auth Service with global state which lives in memory an saved the current logged in user and the service provides an api to get the logged in user
    - Use Auth Service in middlewear to get the logged in user


Problems and Bugs
- fix sqlite wal mode and connection see: https://github.com/KochKevin/effective-spoon-v2/pull/6

Backlog
- Shopping Cart Card größe anpassen
- Use env to point frontend api to backend

- Transactions should save the timestamp when they are created and should be sorted based on this timestamp when the balance is loaded

- auto tools installer task
    - use go get -tool for adding backend tools, see: https://pkg.go.dev/cmd/go#:~:text=The%20%2Dtool%20flag%20instructs%20go%20to%20add%20a%20matching%20tool%20line%20to%20go%2Emod%20for%20each%20listed%20package%2E%20If%20%2Dtool%20is%20used%20with%20%40none%2C%20the%20line%20will%20be%20removed%2E
    - add an task to auto install all go/backend tools and frontend tools (npm install)


- Combine Generate Transaction and Checkout to one Checkout Function which returns an transaction (shoppingcarts/model)


- In main.go, if goose.Up fails, the application logs the error but continues running. It is safer to log the error and terminate the application (log.Fatal) to avoid running with an inconsistent database schema.

- In SaveShoppingCart (backend/internal/shoppingcarts/sqlite/repo.go), the error returned by UpdateShoppingCart is completely ignored. If the update fails, the function will still return nil (success), which can lead to silent data loss or inconsistent state.

- In all handlers, WithTx is called with context.Background() instead of r.Context(). This ignores the request context, meaning database transactions won't be canceled if the client disconnects or the request times out. Please use r.Context().

- Using raw strings like "user_id" as context keys is a Go anti-pattern. Consider defining a custom unexported type for context keys to avoid potential collisions. -> Add an project wide http package with typed context keys and custom middlewear

- Several slog.Error calls pass err directly as the second argument instead of using key-value pairs (e.g., "error", err).


To thinker about
- add more auth to shopping cart
    - check if the user which is the shopping cart owner is the same as the one who is increasing, decreasing or checking it out
