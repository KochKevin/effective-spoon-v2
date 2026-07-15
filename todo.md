- Add User System
    - Shopping Carts can make transactions and withdrawl from users money



- Auth System
    - Provide test endpoint to login with user code
    - Auth Service with global state which lives in memory an saved the current logged in user and the service provides an api to get the logged in user
    - Use Auth Service in middlewear to get the logged in user


Problems and Bugs

- Database locked error when doing to many things

Backlog
- Shopping Cart Card größe anpassen
- Use env to point frontend api to backend

- Transactions should save the timestamp when they are created and should be sorted based on this timestamp when the balance is loaded

- auto tools installer task
    - use go get -tool for adding backend tools, see: https://pkg.go.dev/cmd/go#:~:text=The%20%2Dtool%20flag%20instructs%20go%20to%20add%20a%20matching%20tool%20line%20to%20go%2Emod%20for%20each%20listed%20package%2E%20If%20%2Dtool%20is%20used%20with%20%40none%2C%20the%20line%20will%20be%20removed%2E
    - add an task to auto install all go/backend tools and frontend tools (npm install)


- Combine Generate Transaction and Checkout to one Checkout Function which returns an transaction (shoppingcarts/model)