package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

func compareRequestToResponse(createRequest CreateModelRequest, response ModelResponse) {
	ExpectWithOffset(1, createRequest.TheString).To(Equal(response.TheString))
	ExpectWithOffset(1, createRequest.TheInt).To(Equal(response.TheInt))
	ExpectWithOffset(1, createRequest.ComplexType).To(Equal(response.ComplexType))
}

func compareEmbeddedRequestToResponse(createRequest CreateEmbeddedRequest, response EmbeddedResponse) {
	ExpectWithOffset(1, createRequest.TheString).To(Equal(response.TheString))
	ExpectWithOffset(1, createRequest.TheInt).To(Equal(response.TheInt))
	ExpectWithOffset(1, createRequest.ComplexType).To(Equal(response.ComplexType))
}

func TestEqual(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Equal Suite")
}

var _ = Describe("Equal", func() {
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

	Describe("compare simple field", func() {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.TheString = createModelRequest.TheString + " - not"
		}
		Context("createModelRequest.TheString", func() {
			It("should be createModelRequest2.TheString", func() {
				Expect(createModelRequest.TheString).To(Equal(createModelRequest2.TheString))
			})
		})
	})

	Describe("compare struct", func() {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.TheInt = createModelRequest.TheInt + 5
		}
		Context("createModelRequest", func() {
			It("should be createModelRequest2", func() {
				Expect(createModelRequest).To(Equal(createModelRequest2))
			})
		})
	})

	Describe("compare inner struct", func() {
		createModelRequest2 := cloneCreateModelRequest(createModelRequest)
		if shouldFail {
			createModelRequest2.ComplexType.SubString = createModelRequest.ComplexType.SubString + " - not"
		}
		Context("createModelRequest", func() {
			It("should be createModelRequest2", func() {
				Expect(createModelRequest).To(Equal(createModelRequest2))
			})
		})
	})

	Describe("compare almost equivalent struct", func() {
		response := createModelResponseFromRequest(createModelRequest)
		if shouldFail {
			response.TheString = createModelRequest.TheString + " - not"
		}
		Context("createModelRequest", func() {
			It("should be response", func() {
				compareRequestToResponse(createModelRequest, response)
			})
		})
	})

	Describe("compare embedded model", func() {
		request := CreateEmbeddedRequest{BaseModel: baseModel}
		request2 := CreateEmbeddedRequest{BaseModel: baseModel}
		if shouldFail {
			request2.TheString = request.TheString + " - not"
		}
		Context("request", func() {
			It("should be request2", func() {
				Expect(request).To(Equal(request2))
			})
		})
	})

	Describe("compare almost equivalent embedded model", func() {
		request := CreateEmbeddedRequest{BaseModel: baseModel}
		response := createEmbeddedResponseFromRequest(request)
		if shouldFail {
			response.TheString = request.TheString + " - not"
		}
		Context("request", func() {
			It("should be response", func() {
				compareEmbeddedRequestToResponse(request, response)
			})
		})
	})

	Describe("compare struct with time", func() {
		modelWithTime2 := cloneModelWithTime(modelWithTime)
		if shouldFail {
			modelWithTime2.TheTime = time.Now()
		}
		Context("modelWithTime", func() {
			It("should be modelWithTime2", func() {
				Expect(modelWithTime).To(Equal(modelWithTime2))
			})
		})
	})

	Describe("compare struct with duration", func() {
		modelWithTime2 := cloneModelWithTime(modelWithTime)
		if shouldFail {
			modelWithTime2.TheDuration = 50 * time.Second
		}
		Context("modelWithTime", func() {
			It("should be modelWithTime2", func() {
				Expect(modelWithTime).To(Equal(modelWithTime2))
			})
		})
	})

	Describe("compare small map", func() {
		firstMap := map[string]int{"foo": 10, "bar": 20}
		secondMap := map[string]int{"foo": 10, "bar": 20}
		if shouldFail {
			secondMap["foo"] = 15
		}
		Context("firstMap", func() {
			It("should be secondMap", func() {
				Expect(firstMap).To(Equal(secondMap))
			})
		})
	})

	Describe("compare large map", func() {
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
		Context("firstMap", func() {
			It("should be secondMap", func() {
				Expect(firstMap).To(Equal(secondMap))
			})
		})
	})

	Describe("compare small list", func() {
		firstList := []string{"foo", "bar"}
		secondList := []string{"foo", "bar"}
		if shouldFail {
			secondList[1] = "baz"
		}
		Context("firstList", func() {
			It("should be secondList", func() {
				Expect(firstList).To(Equal(secondList))
			})
		})
	})

	Describe("compare large list", func() {
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
		Context("firstList", func() {
			It("should be secondList", func() {
				Expect(firstList).To(Equal(secondList))
			})
		})
	})

})
