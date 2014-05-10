package router

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

// Params map path keys to values.  For example, if your route has the path pattern:
//  /person/:person_id/pets/:pet_type
// Then a correct Params map would lool like:
//  router.Params{
//    "person_id": "123",
//    "pet_type": "cats",
//  }
type Params map[string]string

// A Route defines properties of an HTTP endpoint.  At runtime, the router will
// associate each Route with a http.Handler object, and use the Route properties
// to determine which Handler should be invoked.
//
// Currently, properties used in match are such as Method and Path.
//
// Method is one of the following:
//  GET PUT POST DELETE
//
// Path conforms to Pat-style pattern matching. The following docs are taken from
// http://godoc.org/github.com/bmizerany/pat#PatternServeMux
//
// Path Patterns may contain literals or captures. Capture names start with a colon
// and consist of letters A-Z, a-z, _, and 0-9. The rest of the pattern
// matches literally. The portion of the URL matching each name ends with an
// occurrence of the character in the pattern immediately following the name,
// or a /, whichever comes first. It is possible for a name to match the empty
// string.
//
// Example pattern with one capture:
//   /hello/:name
// Will match:
//   /hello/blake
//   /hello/keith
// Will not match:
//   /hello/blake/
//   /hello/blake/foo
//   /foo
//   /foo/bar
//
// Example 2:
//    /hello/:name/
// Will match:
//   /hello/blake/
//   /hello/keith/foo
//   /hello/blake
//   /hello/keith
// Will not match:
//   /foo
//   /foo/bar
type Route struct {
	// Handler is a key specifying which HTTP handler the router
	// should associate with the endpoint at runtime.
	Handler string
	// Method is one of the following: GET,PUT,POST,DELETE
	Method string
	// Path contains a path pattern
	Path string
}

// PathWithParams combines the route's path pattern with a Params map
// to produce a valid path.
func (r Route) PathWithParams(params Params) (string, error) {
	components := strings.Split(r.Path, "/")
	for i, c := range components {
		if len(c) == 0 {
			continue
		}
		if c[0] == ':' {
			val, ok := params[c[1:]]
			if !ok {
				return "", fmt.Errorf("missing param %s", c)
			}
			components[i] = val
		}
	}

	u, err := url.Parse(strings.Join(components, "/"))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Routes is a Route collection.
type Routes []Route

// RouteForHandler looks up a Route by it's Handler key.
func (r Routes) RouteForHandler(handler string) (Route, bool) {
	for _, route := range r {
		if route.Handler == handler {
			return route, true
		}
	}
	return Route{}, false
}

// PathForHandler looks up a Route by it's Handler key and computes it's path
// with a given Params map.
func (r Routes) PathForHandler(handler string, params Params) (string, error) {
	route, ok := r.RouteForHandler(handler)
	if !ok {
		return "", fmt.Errorf("No route exists for handler %", handler)
	}
	return route.PathWithParams(params)
}

// Router is deprecated, please use router.NewRouter() instead
func (r Routes) Router(handlers Handlers) (http.Handler, error) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("\n\033[0;35m%s\033[0m%s:%d:%s\n", "WARNING:", file, line, " Routes.Router() is deprecated, please use router.NewRouter() instead")
	return NewRouter(r, handlers)
}
