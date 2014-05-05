package router_test

import (
	"github.com/cloudfoundry/gunk/test_server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/router"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Router", func() {
	Describe("Route", func() {
		var route router.Route

		Describe("PathWithParams", func() {
			BeforeEach(func() {
				route = router.Route{
					Handler: "whatevz",
					Method:  "GET",
					Path:    "/a/path/:param/with/:many_things/:many/in/:it",
				}
			})

			It("should return a url with all :entries populated by the passed in hash", func() {
				Ω(route.PathWithParams(router.Params{
					"param":       "1",
					"many_things": "2",
					"many":        "a space",
					"it":          "4",
				})).Should(Equal(`/a/path/1/with/2/a%20space/in/4`))
			})

			Context("when the hash is missing params", func() {
				It("should error", func() {
					_, err := route.PathWithParams(router.Params{
						"param": "1",
						"many":  "2",
						"it":    "4",
					})
					Ω(err).Should(HaveOccurred())
				})
			})

			Context("when the hash has extra params", func() {
				It("should totally not care", func() {
					Ω(route.PathWithParams(router.Params{
						"param":       "1",
						"many_things": "2",
						"many":        "a space",
						"it":          "4",
						"donut":       "bacon",
					})).Should(Equal(`/a/path/1/with/2/a%20space/in/4`))
				})
			})

			Context("with a trailing slash", func() {
				It("should work", func() {
					route = router.Route{
						Handler: "whatevz",
						Method:  "GET",
						Path:    "/a/path/:param/",
					}
					Ω(route.PathWithParams(router.Params{
						"param": "1",
					})).Should(Equal(`/a/path/1/`))
				})
			})
		})
	})

	Describe("Routes", func() {
		var routes router.Routes

		Describe("RouteForHandler", func() {
			BeforeEach(func() {
				routes = router.Routes{
					{Path: "/something", Method: "GET", Handler: "getter"},
					{Path: "/something", Method: "POST", Handler: "poster"},
					{Path: "/something", Method: "PuT", Handler: "putter"},
					{Path: "/something", Method: "DELETE", Handler: "deleter"},
				}
			})

			Context("when the route is present", func() {
				It("returns the route with the matching handler name", func() {
					route, ok := routes.RouteForHandler("getter")
					Ω(ok).Should(BeTrue())
					Ω(route.Method).Should(Equal("GET"))
				})
			})

			Context("when the route is not present", func() {
				It("returns falseness", func() {
					route, ok := routes.RouteForHandler("orangutanger")
					Ω(ok).Should(BeFalse())
					Ω(route).Should(BeZero())
				})
			})
		})

		Describe("PathForHandler", func() {
			BeforeEach(func() {
				routes = router.Routes{
					{
						Handler: "whatevz",
						Method:  "GET",
						Path:    "/a/path/:param/with/:many_things/:many/in/:it",
					},
				}
			})

			Context("when the route is present", func() {
				It("returns the route with the matching handler name", func() {
					path, err := routes.PathForHandler("whatevz", router.Params{
						"param":       "1",
						"many_things": "2",
						"many":        "a space",
						"it":          "4",
					})
					Expect(err).NotTo(HaveOccurred())
					Ω(path).Should(Equal(`/a/path/1/with/2/a%20space/in/4`))
				})

				Context("when the route is not present", func() {
					It("returns an error", func() {
						_, err := routes.PathForHandler("foo", router.Params{
							"param":       "1",
							"many_things": "2",
							"many":        "a space",
							"it":          "4",
						})
						Expect(err).To(HaveOccurred())
					})
				})

				Context("when the hash is missing params", func() {
					It("should error", func() {
						_, err := routes.PathForHandler("whatevz", router.Params{
							"param": "1",
							"many":  "2",
							"it":    "4",
						})
						Ω(err).Should(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("Router", func() {
		var r http.Handler
		var err error
		var routes = router.Routes{
			{Path: "/something", Method: "GET", Handler: "getter"},
			{Path: "/something", Method: "POST", Handler: "poster"},
			{Path: "/something", Method: "PuT", Handler: "putter"},
			{Path: "/something", Method: "DELETE", Handler: "deleter"},
		}

		Context("when all the handlers are present", func() {
			var resp *httptest.ResponseRecorder
			var handlers = router.Handlers{
				"getter":  test_server.Respond(http.StatusOK, "get response"),
				"poster":  test_server.Respond(http.StatusOK, "post response"),
				"putter":  test_server.Respond(http.StatusOK, "put response"),
				"deleter": test_server.Respond(http.StatusOK, "delete response"),
			}
			BeforeEach(func() {
				resp = httptest.NewRecorder()
				r, err = router.NewRouter(routes, handlers)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("makes GET handlers", func() {
				req, _ := http.NewRequest("GET", "/something", nil)

				r.ServeHTTP(resp, req)
				Ω(resp.Body.String()).Should(Equal("get response"))
			})

			It("makes POST handlers", func() {
				req, _ := http.NewRequest("POST", "/something", nil)

				r.ServeHTTP(resp, req)
				Ω(resp.Body.String()).Should(Equal("post response"))
			})

			It("makes PUT handlers", func() {
				req, _ := http.NewRequest("PUT", "/something", nil)

				r.ServeHTTP(resp, req)
				Ω(resp.Body.String()).Should(Equal("put response"))
			})

			It("makes DELETE handlers", func() {
				req, _ := http.NewRequest("DELETE", "/something", nil)

				r.ServeHTTP(resp, req)
				Ω(resp.Body.String()).Should(Equal("delete response"))
			})
		})

		Context("when a handler is missing", func() {
			var incompleteHandlers = router.Handlers{
				"getter": test_server.Respond(http.StatusOK, "get response"),
			}
			It("should error", func() {
				r, err = router.NewRouter(routes, incompleteHandlers)

				Ω(err).Should(HaveOccurred())
			})
		})

		Context("with an invalid method", func() {
			var invalidRoutes = router.Routes{
				{Path: "/something", Method: "SMELL", Handler: "smeller"},
			}

			It("should error", func() {
				handlers := router.Handlers{
					"smeller": test_server.Respond(http.StatusOK, "smelt response"),
				}
				r, err = router.NewRouter(invalidRoutes, handlers)

				Ω(err).Should(HaveOccurred())
			})
		})
	})
})
