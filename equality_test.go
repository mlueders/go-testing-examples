package main

import (
	"fmt"
	. "github.com/franela/goblin"
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
	TheTime     time.Time
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
		TheTime:     modelWithTime.TheTime,
		TheDuration: modelWithTime.TheDuration,
	}
}

func compareRequestToResponse(g *G, createRequest CreateModelRequest, response ModelResponse) {
	g.Assert(createRequest.TheString).Equal(response.TheString)
	g.Assert(createRequest.TheInt).Equal(response.TheInt)
	g.Assert(createRequest.ComplexType).Equal(response.ComplexType)
}

func compareEmbeddedRequestToResponse(g *G, createRequest CreateEmbeddedRequest, response EmbeddedResponse) {
	g.Assert(createRequest.TheString).Equal(response.TheString)
	g.Assert(createRequest.TheInt).Equal(response.TheInt)
	g.Assert(createRequest.ComplexType).Equal(response.ComplexType)
}

func TestEqual(t *testing.T) {
	g := Goblin(t)

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
		TheTime:     time.Now(),
		TheDuration: 10 * time.Second,
	}

	g.Describe("Equal", func() {
		g.It("should compare simple field", func() {
			createModelRequest2 := cloneCreateModelRequest(createModelRequest)
			if shouldFail {
				createModelRequest2.TheString = createModelRequest.TheString + " - not"
			}
			g.Assert(createModelRequest).Equal(createModelRequest2)
		})

		g.It("should compare struct", func() {
			createModelRequest2 := cloneCreateModelRequest(createModelRequest)
			if shouldFail {
				createModelRequest2.TheInt = createModelRequest.TheInt + 5
			}
			g.Assert(createModelRequest).Equal(createModelRequest2)
		})

		g.It("should compare inner struct", func() {
			createModelRequest2 := cloneCreateModelRequest(createModelRequest)
			if shouldFail {
				createModelRequest2.ComplexType.SubString = createModelRequest.ComplexType.SubString + " - not"
			}
			g.Assert(createModelRequest).Equal(createModelRequest2)
		})

		g.It("should compare almost equivalent struct", func() {
			response := createModelResponseFromRequest(createModelRequest)
			if shouldFail {
				response.TheString = createModelRequest.TheString + " - not"
			}
			compareRequestToResponse(g, createModelRequest, response)
		})

		g.It("should compare embedded model", func() {
			request := CreateEmbeddedRequest{BaseModel: baseModel}
			request2 := CreateEmbeddedRequest{BaseModel: baseModel}
			if shouldFail {
				request2.TheString = request.TheString + " - not"
			}
			g.Assert(request).Equal(request2)
		})

		g.It("should compare almost equivalent embedded model", func() {
			request := CreateEmbeddedRequest{BaseModel: baseModel}
			response := createEmbeddedResponseFromRequest(request)
			if shouldFail {
				response.TheString = request.TheString + " - not"
			}
			compareEmbeddedRequestToResponse(g, request, response)
		})

		g.It("should compare struct with time", func() {
			modelWithTime2 := cloneModelWithTime(modelWithTime)
			if shouldFail {
				modelWithTime2.TheTime = time.Now()
			}
			g.Assert(modelWithTime).Equal(modelWithTime2)
		})

		g.It("should compare struct with duration", func() {
			modelWithTime2 := cloneModelWithTime(modelWithTime)
			if shouldFail {
				modelWithTime2.TheDuration = 50 * time.Second
			}
			g.Assert(modelWithTime).Equal(modelWithTime2)
		})

		g.It("should compare small map", func() {
			firstMap := map[string]int{"foo": 10, "bar": 20}
			secondMap := map[string]int{"foo": 10, "bar": 20}
			if shouldFail {
				secondMap["foo"] = 15
			}
			g.Assert(firstMap).Equal(secondMap)
		})

		g.It("should compare large map", func() {
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
			g.Assert(firstMap).Equal(secondMap)
		})

		g.It("should compare small list", func() {
			firstList := []string{"foo", "bar"}
			secondList := []string{"foo", "bar"}
			if shouldFail {
				secondList[1] = "baz"
			}
			g.Assert(firstList).Equal(secondList)
		})

		g.It("should compare large list", func() {
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
			g.Assert(firstList).Equal(secondList)
		})
	})

}
