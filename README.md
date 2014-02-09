martini-experiments
===================

Playing with martini

Adding support for CORS Access-Control-Allow-Methods headers and Allow headers on a 405 are easy in my old code. With Martini there is no support for getting the methods. 

First attempt: used .All() on a route to catch unknown methods and provide either a 405 or CORS on OPTION. This had the disadvantage of having to enter the allowed methods for each path.

Second attempt: add support for Route.MethodsFor (PR sent) and use that in a NotFound handler to provide both CORS and 405 responses.

Third attempt: Handling CORS in two places (middleware and NotFound) seems silly but Routes is only available within the Router which happens in Action after the middleware in Martini Classic. For the second attempt I made it available for NotFound handlers as well. This time I have added it into the ClassicMartini function so it is available for middlware as well.