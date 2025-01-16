package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/event"
	"github.com/guilehm/echo-vision/internal/infra/web"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = XDescribe("Event Handler", func() {
	Context("Create Event", func() {
		It("should create a event successfully", func() {
			// Arrange
			input := ports.EventCreateInput{
				EventType: domain.EventTypeImageAnalysis.String(),
				SubType:   domain.EventSubTypeDetectLabels.String(),
			}
			b, _ := json.Marshal(input)

			// Act
			resp, err := http.Post(fmt.Sprintf("%s/events", server.URL), "application/json", bytes.NewReader(b))
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var apiResp web.ApiResponse[ports.EventCreateResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Data.ID).ToNot(BeNil())
			Expect(apiResp.Error).To(BeEmpty())

			createdEvent := entClient.Event.Query().
				Where(event.IDEQ(apiResp.Data.ID)).
				OnlyX(ctx)

			Expect(createdEvent.Type).To(Equal(domain.EventTypeImageAnalysis))
			Expect(createdEvent.SubType).To(Equal(domain.EventSubTypeDetectLabels))
		})
	})
})
