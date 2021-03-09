package main

import (
	"encoding/json"
	"errors"
	"github.com/CodeWithSameera/Vehicles/handlers"
	"github.com/CodeWithSameera/Vehicles/helpers"
	"github.com/CodeWithSameera/Vehicles/model"
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)
type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func main() {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := helpers.GoDotEnvVariable("AUDIENCE")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, true)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := helpers.GoDotEnvVariable("DOMAIN")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})



	r.Handle("/api/vehicles", jwtMiddleware.Handler(handlers.GetAllVehiclesHandler)).Methods("GET")
	r.Handle("/api/vehicles/{vno}", jwtMiddleware.Handler(handlers.SearchVehiclesHandler)).Methods("GET")
	r.Handle("/api/vehicles", jwtMiddleware.Handler(handlers.CreateVehiclesHandler)).Methods("POST")
	r.Handle("/api/vehicles/{vno}", jwtMiddleware.Handler(handlers.UpdateVehiclesHandler)).Methods("PUT")
	r.Handle("/api/vehicles/{vno}", jwtMiddleware.Handler(handlers.DeleteVehiclesHandler)).Methods("DELETE")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST","PUT","DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})
	port := helpers.GoDotEnvVariable("PORT")
	db :=helpers.ConnectDB();
	model.DBMigrate(db);
	http.ListenAndServe(port, corsWrapper.Handler(r))
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(helpers.GoDotEnvVariable("DOMAIN")+".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}