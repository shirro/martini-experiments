martini-experiments
===================

Playing with martini


My main interest in Martini is the DI as I am also doing Angular frontend and my head is in that space at the moment.

My concern is the routing. I have a previous homebrew arrangement built on net/http where I register an interface{} and use reflection to bind appropriately named methods to http methods and the routes are basically map[string]interface{}.

As a result CORS Access-Control-Allow-Methods headers and Allow headers on a 405 are easy in my old code. With Martini, short of scanning the entire route slice to find matching methods I have to provide this information manually to a catchall .All() handler for a route.
