package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

var r *chi.Mux

func RegisterRoutes(handler *handler) *chi.Mux {
	r = chi.NewRouter()

	r.Route("/products", func(r chi.Router) {
		r.Post("/", handler.createProduct)
		r.Get("/", handler.listProducts)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.getProduct)
			r.Patch("/", handler.updateProduct)
			r.Delete("/", handler.deleteProduct)
		})
	})

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", handler.createOrder)
		r.Get("/", handler.listOrders)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.getOrder)
			r.Delete("/", handler.deleteOrder)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", handler.createUser)
		r.Get("/", handler.listUsers)
		r.Patch("/", handler.updateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", handler.deleteUser)
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", handler.loginUser)
		})

		r.Route("/logout", func(r chi.Router) {
			r.Post("/", handler.logoutUser)
		})
	})

	r.Route("/tokens", func(r chi.Router) {
		r.Route("/renew", func(r chi.Router) {
			r.Post("/", handler.renewAccessToken)
		})

		r.Route("/revoke/{id}", func(r chi.Router) {
			r.Post("/", handler.revokeSession)
		})
	})

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, r)
}
