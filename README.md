martini-experiments
===================

Playing with martini

Adding support for CORS Access-Control-Allow-Methods headers and Allow headers on a 405 are easy in my old code. With Martini there is no support for getting the methods. 

My first attempt used .All() on a route to catch unknown methods and provide either a 405 or CORS on OPTION. This had the disadvantage of having to enter the allowed methods for each path.

My second attempt was to add support for Route.MethodsFor (PR sent) and use that in a NotFound handler to provide both CORS and 405 responses.