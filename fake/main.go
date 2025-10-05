package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jaswdr/faker"
)

type Request struct {
	Function string                 `json:"function"`
	Params   map[string]any `json:"params"`
	Context  map[string]any `json:"context"`
}

type Response struct {
	Value any `json:"value"`
	Type  string      `json:"type"`
	Error string      `json:"error,omitempty"`
}

var fake faker.Faker

func main() {
	// Initialize faker with seed from params or use timestamp
	fake = faker.New()

	var req Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		sendError(fmt.Sprintf("Invalid request: %v", err))
		return
	}

	// Allow seeding for reproducible data
	if seed, ok := req.Params["seed"].(float64); ok {
		_ = seed // Seeding support for future enhancement
		// fake = faker.NewWithSeed(...)
	}

	var result any
	var resultType string
	var err error

	switch req.Function {
	// Person functions
	case "name":
		result, resultType, err = handleName(req.Params)
	case "firstName":
		result, resultType, err = handleFirstName(req.Params)
	case "lastName":
		result, resultType, err = handleLastName(req.Params)
	case "email":
		result, resultType, err = handleEmail(req.Params)
	case "phone":
		result, resultType, err = handlePhone(req.Params)
	case "username":
		result, resultType, err = handleUsername(req.Params)

	// Internet functions
	case "url":
		result, resultType, err = handleURL(req.Params)
	case "domain":
		result, resultType, err = handleDomain(req.Params)
	case "ipv4":
		result, resultType, err = handleIPv4(req.Params)
	case "ipv6":
		result, resultType, err = handleIPv6(req.Params)
	case "userAgent":
		result, resultType, err = handleUserAgent(req.Params)

	// Company functions
	case "company":
		result, resultType, err = handleCompany(req.Params)
	case "jobTitle":
		result, resultType, err = handleJobTitle(req.Params)

	// Address functions
	case "address":
		result, resultType, err = handleAddress(req.Params)
	case "city":
		result, resultType, err = handleCity(req.Params)
	case "state":
		result, resultType, err = handleState(req.Params)
	case "country":
		result, resultType, err = handleCountry(req.Params)
	case "zipCode":
		result, resultType, err = handleZipCode(req.Params)
	case "latitude":
		result, resultType, err = handleLatitude(req.Params)
	case "longitude":
		result, resultType, err = handleLongitude(req.Params)

	// Text functions
	case "word":
		result, resultType, err = handleWord(req.Params)
	case "sentence":
		result, resultType, err = handleSentence(req.Params)
	case "paragraph":
		result, resultType, err = handleParagraph(req.Params)
	case "lorem":
		result, resultType, err = handleLorem(req.Params)

	// Commerce functions
	case "product":
		result, resultType, err = handleProduct(req.Params)
	case "price":
		result, resultType, err = handlePrice(req.Params)
	case "currency":
		result, resultType, err = handleCurrency(req.Params)

	// Color functions
	case "color":
		result, resultType, err = handleColor(req.Params)
	case "hexColor":
		result, resultType, err = handleHexColor(req.Params)

	// Misc functions
	case "creditCard":
		result, resultType, err = handleCreditCard(req.Params)

	default:
		sendError(fmt.Sprintf("Unknown function: %s", req.Function))
		return
	}

	if err != nil {
		sendError(err.Error())
		return
	}

	sendResponse(result, resultType)
}

// Person functions

func handleName(params map[string]any) (any, string, error) {
	return fake.Person().Name(), "string", nil
}

func handleFirstName(params map[string]any) (any, string, error) {
	return fake.Person().FirstName(), "string", nil
}

func handleLastName(params map[string]any) (any, string, error) {
	return fake.Person().LastName(), "string", nil
}

func handleEmail(params map[string]any) (any, string, error) {
	return fake.Internet().Email(), "string", nil
}

func handlePhone(params map[string]any) (any, string, error) {
	return fake.Phone().Number(), "string", nil
}

func handleUsername(params map[string]any) (any, string, error) {
	return fake.Internet().User(), "string", nil
}

// Internet functions

func handleURL(params map[string]any) (any, string, error) {
	return fake.Internet().URL(), "string", nil
}

func handleDomain(params map[string]any) (any, string, error) {
	return fake.Internet().Domain(), "string", nil
}

func handleIPv4(params map[string]any) (any, string, error) {
	return fake.Internet().Ipv4(), "string", nil
}

func handleIPv6(params map[string]any) (any, string, error) {
	return fake.Internet().Ipv6(), "string", nil
}

func handleUserAgent(params map[string]any) (any, string, error) {
	return fake.UserAgent().UserAgent(), "string", nil
}

// Company functions

func handleCompany(params map[string]any) (any, string, error) {
	return fake.Company().Name(), "string", nil
}

func handleJobTitle(params map[string]any) (any, string, error) {
	return fake.Company().JobTitle(), "string", nil
}

// Address functions

func handleAddress(params map[string]any) (any, string, error) {
	return fake.Address().Address(), "string", nil
}

func handleCity(params map[string]any) (any, string, error) {
	return fake.Address().City(), "string", nil
}

func handleState(params map[string]any) (any, string, error) {
	return fake.Address().State(), "string", nil
}

func handleCountry(params map[string]any) (any, string, error) {
	return fake.Address().Country(), "string", nil
}

func handleZipCode(params map[string]any) (any, string, error) {
	return fake.Address().PostCode(), "string", nil
}

func handleLatitude(params map[string]any) (any, string, error) {
	return fake.Address().Latitude(), "float", nil
}

func handleLongitude(params map[string]any) (any, string, error) {
	return fake.Address().Longitude(), "float", nil
}

// Text functions

func handleWord(params map[string]any) (any, string, error) {
	return fake.Lorem().Word(), "string", nil
}

func handleSentence(params map[string]any) (any, string, error) {
	words := getInt(params, "words", 10)
	return fake.Lorem().Sentence(words), "string", nil
}

func handleParagraph(params map[string]any) (any, string, error) {
	sentences := getInt(params, "sentences", 3)
	return fake.Lorem().Paragraph(sentences), "string", nil
}

func handleLorem(params map[string]any) (any, string, error) {
	words := getInt(params, "words", 50)
	return fake.Lorem().Text(words), "string", nil
}

// Commerce functions

func handleProduct(params map[string]any) (any, string, error) {
	return fake.Beer().Name(), "string", nil
}

func handlePrice(params map[string]any) (any, string, error) {
	min := getFloat(params, "min", 1.0)
	max := getFloat(params, "max", 1000.0)
	// Generate random price between min and max with 2 decimal places
	price := fake.Float64(2, int(min), int(max))
	return price, "float", nil
}

func handleCurrency(params map[string]any) (any, string, error) {
	return fake.Currency().Currency(), "string", nil
}

// Color functions

func handleColor(params map[string]any) (any, string, error) {
	return fake.Color().ColorName(), "string", nil
}

func handleHexColor(params map[string]any) (any, string, error) {
	return fake.Color().Hex(), "string", nil
}

// Misc functions

func handleCreditCard(params map[string]any) (any, string, error) {
	ccType := getString(params, "type", "")

	switch ccType {
	case "visa":
		return fake.Payment().CreditCardNumber(), "string", nil
	case "mastercard":
		return fake.Payment().CreditCardNumber(), "string", nil
	default:
		return fake.Payment().CreditCardNumber(), "string", nil
	}
}

// Helper functions

func getString(params map[string]any, key, defaultValue string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

func getInt(params map[string]any, key string, defaultValue int) int {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return defaultValue
}

func getFloat(params map[string]any, key string, defaultValue float64) float64 {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		}
	}
	return defaultValue
}

func sendResponse(value any, valueType string) {
	resp := Response{
		Value: value,
		Type:  valueType,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}

func sendError(message string) {
	resp := Response{
		Error: message,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
	os.Exit(1)
}

