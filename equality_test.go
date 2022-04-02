package main

import (
	"fmt"
	"github.com/go-test/deep"
	"os"
	"testing"
	"time"
)

type CreateModelRequest struct {
	TheString   string
	TheInt      int
	ComplexType SubModel
}

type ModelResponse struct {
	ID          string
	TheString   string
	TheInt      int
	ComplexType SubModel
}

type SubModel struct {
	SubString string
	SubFloat  float64
}

type BaseModel struct {
	TheString   string
	TheInt      int
	ComplexType SubModel
}

type CreateEmbeddedRequest struct {
	BaseModel
}

type EmbeddedResponse struct {
	BaseModel
	ID string
}

type ModelWithTime struct {
	TheTime time.Time
	TheDuration time.Duration
}

func cloneCreateModelRequest(request CreateModelRequest) CreateModelRequest {
	return CreateModelRequest{
		TheString:   request.TheString,
		TheInt:      request.TheInt,
		ComplexType: request.ComplexType,
	}
}

func createModelResponseFromRequest(request CreateModelRequest) ModelResponse {
	return ModelResponse{
		ID:          "someid",
		TheString:   request.TheString,
		TheInt:      request.TheInt,
		ComplexType: request.ComplexType,
	}
}

func createEmbeddedResponseFromRequest(request CreateEmbeddedRequest) EmbeddedResponse {
	return EmbeddedResponse{
		ID:        "someid",
		BaseModel: request.BaseModel,
	}
}

func cloneModelWithTime(modelWithTime ModelWithTime) ModelWithTime {
	return ModelWithTime{
		TheTime: modelWithTime.TheTime,
		TheDuration: modelWithTime.TheDuration,
	}
}

func compareRequestToResponse(t *testing.T, createRequest CreateModelRequest, response ModelResponse) {
	assertEqual(t, createRequest.TheString, response.TheString)
	assertEqual(t, createRequest.TheInt, response.TheInt)
	assertEqual(t, createRequest.ComplexType, response.ComplexType)
}

func compareEmbeddedRequestToResponse(t *testing.T, createRequest CreateEmbeddedRequest, response EmbeddedResponse) {
	assertEqual(t, createRequest.TheString, response.TheString)
	assertEqual(t, createRequest.TheInt, response.TheInt)
	assertEqual(t, createRequest.ComplexType, response.ComplexType)
}

func assertEqual(t *testing.T, a, b interface{}) {
	if diff := deep.Equal(a, b); diff != nil {
		t.Error(diff)
	}
}

func TestEqual(t *testing.T) {
	shouldFail := os.Getenv("SHOULD_FAIL") == "true"
	createModelRequest := CreateModelRequest{
		TheString: "some string",
		TheInt:    123,
		ComplexType: SubModel{
			SubString: "sub string",
			SubFloat:  123.456,
		},
	}
	baseModel := BaseModel{
		TheString: "some string",
		TheInt:    123,
		ComplexType: SubModel{
			SubString: "sub string",
			SubFloat:  123.456,
		},
	}
	modelWithTime := ModelWithTime{
		TheTime: time.Now(),
		TheDuration: 10 * time.Second,
	}

	t.Run("should compare simple field", func(t *testing.T) {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.TheString = createModelRequest.TheString + " - not"
		}
		assertEqual(t, createModelRequest.TheString, createModelRequest2.TheString)
	})

	t.Run("should compare struct", func(t *testing.T) {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.TheInt = createModelRequest.TheInt + 5
		}
		assertEqual(t, createModelRequest, createModelRequest2)
	})

	t.Run("should compare inner struct", func(t *testing.T) {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.ComplexType.SubString = createModelRequest.ComplexType.SubString + " - not"
		}
		assertEqual(t, createModelRequest, createModelRequest2)
	})

	t.Run("should compare almost equivalent struct", func(t *testing.T) {
		response := createModelResponseFromRequest(createModelRequest)
		if shouldFail {
			response.TheString = createModelRequest.TheString + " - not"
		}
		compareRequestToResponse(t, createModelRequest, response)
	})

	t.Run("should compare embedded model", func(t *testing.T) {
		request := CreateEmbeddedRequest{BaseModel: baseModel}
		request2 := CreateEmbeddedRequest{BaseModel: baseModel}
		if shouldFail {
			request2.TheString = request.TheString + " - not"
		}
		assertEqual(t, request, request2)
	})

	t.Run("should compare almost equivalent embedded model", func(t *testing.T) {
		request := CreateEmbeddedRequest{BaseModel: baseModel}
		response := createEmbeddedResponseFromRequest(request)
		if shouldFail {
			response.TheString = request.TheString + " - not"
		}
		compareEmbeddedRequestToResponse(t, request, response)
	})

	t.Run("should compare struct with time", func(t *testing.T) {
		modelWithTime2 := cloneModelWithTime(modelWithTime)
		if shouldFail {
			modelWithTime2.TheTime = time.Now()
		}
		assertEqual(t, modelWithTime, modelWithTime2)
	})

	t.Run("should compare struct with duration", func(t *testing.T) {
		modelWithTime2 := cloneModelWithTime(modelWithTime)
		if shouldFail {
			modelWithTime2.TheDuration = 50 * time.Second
		}
		assertEqual(t, modelWithTime, modelWithTime2)
	})

	t.Run("should compare small map", func(t *testing.T) {
		firstMap := map[string]int{"foo":10,"bar":20}
		secondMap := map[string]int{"foo":10,"bar":20}
		if shouldFail {
			secondMap["foo"] = 15
		}
		assertEqual(t, firstMap, secondMap)
	})

	t.Run("should compare large map", func(t *testing.T) {
		firstMap := map[string]int{}
		secondMap := map[string]int{}
		for i := 1; i <= 100; i++ {
			firstMap[fmt.Sprintf("%s.%d", "item", i)] = i
			secondMap[fmt.Sprintf("%s.%d", "item", i)] = i
		}
		if shouldFail {
			secondMap["item.1"] = -1
			secondMap["item.50"] = -1
			secondMap["item.99"] = -1
		}
		assertEqual(t, firstMap, secondMap)
	})

	t.Run("should compare small list", func(t *testing.T) {
		firstList := []string{"foo", "bar"}
		secondList := []string{"foo", "bar"}
		if shouldFail {
			secondList[1] = "baz"
		}
		assertEqual(t, firstList, secondList)
	})

	t.Run("should compare large list", func(t *testing.T) {
		var firstList []string
		var secondList []string
		for i := 1; i <= 100; i++ {
			firstList = append(firstList, fmt.Sprintf("%s.%d", "item", i))
			secondList = append(secondList, fmt.Sprintf("%s.%d", "item", i))
		}
		if shouldFail {
			secondList[0] = "item.-1"
			secondList[50] = "item.-1"
			secondList[99] = "item.-1"
		}
		assertEqual(t, firstList, secondList)
	})

}
