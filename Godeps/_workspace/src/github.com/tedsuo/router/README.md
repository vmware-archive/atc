# Router
A router with Pat-style path patterns.

API Docs: https://godoc.org/github.com/tedsuo/router

A key advantage is the ability to creates Routes as structs, indepentent of the router and the corresponding route handlers.  This allows you to reuse the routes in multiple contexts. For example, a proxy server and a backend server can be created by having one set of Routes, but two sets of Handlers (one handler that proxies, another that serves the request). By sharing the same route structure, the type system ensures the two servers stay in sync.  Likewise, your client code can use the routes to construct requests to the server, ensuring the client and server don't drift apart.