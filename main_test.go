package main

import (
	"fmt"
	"os"
	"reflect"
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
	if createRequest.TheString != response.TheString {
		t.Errorf("TheString == %q, want %q", response.TheString, createRequest.TheString)
	}
	if createRequest.TheInt != response.TheInt {
		t.Errorf("TheInt == %q, want %q", response.TheInt, createRequest.TheInt)
	}
	if reflect.DeepEqual(createRequest.ComplexType, response.ComplexType) == false {
		t.Errorf("ComplexType == %v, want %v", createRequest.ComplexType, response.ComplexType)
	}
}

func compareEmbeddedRequestToResponse(t *testing.T, createRequest CreateEmbeddedRequest, response EmbeddedResponse) {
	if createRequest.TheString != response.TheString {
		t.Errorf("TheString == %q, want %q", response.TheString, createRequest.TheString)
	}
	if createRequest.TheInt != response.TheInt {
		t.Errorf("TheInt == %q, want %q", response.TheInt, createRequest.TheInt)
	}
	if reflect.DeepEqual(createRequest.ComplexType, response.ComplexType) == false {
		t.Errorf("ComplexType == %v, want %v", createRequest.ComplexType, response.ComplexType)
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
		if createModelRequest.TheString != createModelRequest2.TheString {
			t.Errorf("TheString == %q, want %q", createModelRequest2.TheString, createModelRequest.TheString)
		}
	})

	t.Run("should compare struct", func(t *testing.T) {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.TheInt = createModelRequest.TheInt + 5
		}
		if reflect.DeepEqual(createModelRequest, createModelRequest2) == false {
			t.Errorf("CreateModelRequest == %v, want %v", createModelRequest2, createModelRequest)
		}
	})

	t.Run("should compare inner struct", func(t *testing.T) {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.ComplexType.SubString = createModelRequest.ComplexType.SubString + " - not"
		}
		if reflect.DeepEqual(createModelRequest, createModelRequest2) == false {
			t.Errorf("CreateModelRequest == %v, want %v", createModelRequest2, createModelRequest)
		}
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
		if reflect.DeepEqual(request, request2) == false {
			t.Errorf("EmbeddedResponse == %v, want %v", request, request2)
		}
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
		if reflect.DeepEqual(modelWithTime, modelWithTime2) == false {
			t.Errorf("ModelWithTime == %v, want %v", modelWithTime, modelWithTime2)
		}
	})

	t.Run("should compare struct with duration", func(t *testing.T) {
		modelWithTime2 := cloneModelWithTime(modelWithTime)
		if shouldFail {
			modelWithTime2.TheDuration = 50 * time.Second
		}
		if reflect.DeepEqual(modelWithTime, modelWithTime2) == false {
			t.Errorf("ModelWithTime == %v, want %v", modelWithTime, modelWithTime2)
		}
	})

	t.Run("should compare small map", func(t *testing.T) {
		firstMap := map[string]int{"foo":10,"bar":20}
		secondMap := map[string]int{"foo":10,"bar":20}
		if shouldFail {
			secondMap["foo"] = 15
		}
		if reflect.DeepEqual(firstMap, secondMap) == false {
			t.Errorf("Small Map == %v, want %v", firstMap, secondMap)
		}
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
		if reflect.DeepEqual(firstMap, secondMap) == false {
			t.Errorf("Large Map == %v, want %v", firstMap, secondMap)
		}
	})

	t.Run("should compare small list", func(t *testing.T) {
		firstList := []string{"foo", "bar"}
		secondList := []string{"foo", "bar"}
		if shouldFail {
			secondList[1] = "baz"
		}
		if reflect.DeepEqual(firstList, secondList) == false {
			t.Errorf("Small List == %v, want %v", firstList, secondList)
		}
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
		if reflect.DeepEqual(firstList, secondList) == false {
			t.Errorf("Large List == %v, want %v", firstList, secondList)
		}
	})

}
